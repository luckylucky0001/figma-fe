"""figmafe match 准确率：以人工确认的 design/component-map.yaml 为标准答案，
看 match（不读映射文件、只靠资产扫描）推荐的首选组件是否一致；缺口（gaps）应判为「需新建」。
用法：python3 experiments/match-accuracy.py"""
import json, subprocess, tempfile, os, yaml

root = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
m = yaml.safe_load(open(os.path.join(root, 'design/component-map.yaml')))
out = tempfile.mkdtemp()
for pf, extra in [('h5', []), ('android', ['--target', 'compose'])]:
    subprocess.run([os.path.join(root, 'bin/figmafe'), 'match', '--figma', 'design/figma-export/order-detail', '.', '--platform', pf, '-q', '-o', out] + extra, cwd=root, check=True, capture_output=True)
    truth = {e['figma']: (e.get(pf) or {}).get('component') for e in m['components']}
    r = json.load(open(os.path.join(out, pf + '.json')))
    ok = 0
    for it in r['items']:
        got = it['candidates'][0]['name'] if it['candidates'] and it['verdict'] != 'new' else None
        hit = got == truth.get(it['figma'])
        ok += hit
        print(f"  [{pf}] {it['figma']:12s} 标准答案={str(truth.get(it['figma'])):12s} match={str(got):14s} {'✓' if hit else '✗'}")
    print(f"{pf}: {ok}/{len(r['items'])}")
