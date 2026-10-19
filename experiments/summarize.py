"""汇总 experiments/runs/*/report.json → summary.md / summary.json"""
import json, glob, os, re
rows = []
for p in sorted(glob.glob(os.path.join(os.path.dirname(__file__), 'runs', '*', 'report.json'))):
    run = os.path.basename(os.path.dirname(p))
    r = json.load(open(p))
    m = r['metrics']
    src_dir = os.path.join(os.path.dirname(p), 'src')
    src = ''
    for f in glob.glob(src_dir + '/**/*', recursive=True):
        if os.path.isfile(f):
            src += open(f, encoding='utf-8', errors='ignore').read()
    changed = open(os.path.join(os.path.dirname(p), 'changed-files.txt')).read()
    extra = [l.split()[-1] for l in changed.splitlines() if l.startswith('??') and 'OrderDetail' not in l]
    countdown_custom = bool(re.search(r'(function|fun|const)\s+\w*Count[Dd]own', src) or re.search(r'setInterval|setTimeout|LaunchedEffect|delay\(', src))
    rows.append(dict(run=run, group=run.split('-')[0], platform=r['platform'], pass_=r['pass'],
        reuse=m['reuseRate'], runtime=m['runtimeReuse'], violations=m['violations'], warnings=m['warnings'],
        pixel=m['pixelMismatch'], nodeErr=m['nodeErrors'], nodeWarn=m['nodeWarnings'], textCov=m['textCoverage'],
        figmaIds=m['figmaIds'], selfBuiltCountdown=countdown_custom, extraFiles=extra, failed=r['failed']))
json.dump(rows, open(os.path.join(os.path.dirname(__file__), 'summary.json'), 'w'), ensure_ascii=False, indent=2)
hdr = '| run | 平台 | 结果 | 复用率(静态/运行时) | 违规 | 警告 | 像素差异 | 节点偏差(严重/轻微) | 文本覆盖 | figma id | 自造 CountDown |\n|---|---|---|---|---|---|---|---|---|---|---|\n'
lines = [f"| {r['run']} | {r['platform']} | {'PASS' if r['pass_'] else 'FAIL'} | {r['reuse']:.0%} / {r['runtime']:.0%} | {r['violations']} | {r['warnings']} | {r['pixel']:.2%} | {r['nodeErr']} / {r['nodeWarn']} | {r['textCov']:.0%} | {r['figmaIds']} | {'是' if r['selfBuiltCountdown'] else '否'} |" for r in rows]
open(os.path.join(os.path.dirname(__file__), 'summary.md'), 'w').write('# 实验汇总（由 summarize.py 生成）\n\n' + hdr + '\n'.join(lines) + '\n')
print(hdr + '\n'.join(lines))
