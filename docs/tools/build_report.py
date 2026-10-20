"""生成 docs/report.html（10 页 A4），再由 shoot.mjs 输出 docs/report.pdf：

    python3 docs/tools/build_report.py && node docs/tools/shoot.mjs html docs/report.html docs/report.pdf

数字全部来自实验数据，不手填：
  experiments/real-match-eval.json / .md   复用清单准确率（real-match-eval.py）
  experiments/real-project/summary.json    真实项目 A / C 对照（real-summary.py）
  experiments/real-project/cost.md         真实项目实验的 token / 耗时
  experiments/summary.json                 demo 仓库新页面实验（summarize.py）
  docs/real-test/large/<系统>/             三个大型系统的体检和复用清单（figmafe reuse / match 的输出）
只有「第一次跑」的结果和对应修正是文字：它们来自修正前的版本，记录在 experiments/match-judgement.md。"""
import json, os, re

root = os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..'))
ex = lambda *p: os.path.join(root, 'experiments', *p)
img = 'tutorial/img/'
rimg = 'real-test/img/'

# ---- 复用清单准确率 ----
ev = json.load(open(ex('real-match-eval.json')))
stack, heldout = {}, set()
for line in open(ex('real-match-eval.md'), encoding='utf-8'):
    m = re.match(r'\| ([\w-]+)(（留出）)? \| ([^|]+) \|', line)
    if m and m.group(1) != '项目':
        stack[m.group(1)] = m.group(3).strip()
        if m.group(2):
            heldout.add(m.group(1))
projects = list(dict.fromkeys(k.split('/')[0] for k in ev))
pages = [('order-detail', '订单详情'), ('holdings', '我的持仓'), ('trade-order', '买入下单'), ('product-detail', '商品详情')]
tot = {s: sum(v[s] for v in ev.values()) for s in '✓△✗'}
n_items = sum(tot.values())
sc_projects = [p for p in projects if p + '/product-detail' in ev]


def cell(v):
    if not v:
        return '<span class="dim">—</span>'
    s = f"✓{v['✓']}"
    if v['△']:
        s += f" <span class='warn'>△{v['△']}</span>"
    if v['✗']:
        s += f" <b class='bad'>✗{v['✗']}</b>"
    return s


def eval_rows(part):
    half = (len(projects) + 1) // 2
    out = []
    for p in projects[:half] if part == 0 else projects[half:]:
        tag = " <span class='tagho'>留出</span>" if p in heldout else ''
        out.append(f"<tr><td>{p}{tag}</td>" + ''.join(f"<td>{cell(ev.get(p + '/' + k))}</td>" for k, _ in pages) + '</tr>')
    return '\n'.join(out)


# ---- 真实项目 A / C 对照 ----
rp = json.load(open(ex('real-project', 'summary.json')))
INVALID = {'Rlm-A-1'}  # A 组 agent 自己找到并运行了 Skill，不能当基线（见 real-project/summary.md）
runs = [r for r in rp.values() if r['run'] not in INVALID]
PNAME = {'Rnb': ('newbee-mall', 'Vue 3 + Vant'), 'Rlm': ('litemall（留出）', 'Vue 2 + Vant 2'),
         'Rqm': ('QMUI_Android', 'XML + QMUI'), 'Rnp': ('NewPipe（留出）', 'XML + AppCompat'),
         'Rnia': ('Now in Android', 'Compose + 自有设计系统')}
grp = lambda rs, g: [r for r in rs if r['group'] == g]
passes = lambda rs: f"{sum(r['pass'] for r in rs)}/{len(rs)}" if rs else '—'
bad = lambda r: r['reinvented'] + r['reinventedWarn'] + r['missed'] + r['handBuiltError'] + r['handBuiltWarn']


def exp_rows():
    out = []
    for p, (name, st) in PNAME.items():
        rs = [r for r in runs if r['project'] == p]
        a, c = grp(rs, 'A') + grp(rs, 'Ahk'), grp(rs, 'C') + grp(rs, 'Chk')
        f = lambda xs, k: ' / '.join(str(k(r)) for r in xs)
        out.append(f"<tr><td>{name}</td><td>{st}</td><td class='bad'>{passes(a)}</td><td>{f(a, lambda r: r['used'])}</td><td>{f(a, bad)}</td>"
                   f"<td class='ok'>{passes(c)}</td><td>{f(c, lambda r: r['used'])}</td><td>{f(c, bad)}</td></tr>")
    return '\n'.join(out)


A_all = [r for r in runs if r['group'].startswith('A')]
C_all = [r for r in runs if r['group'].startswith('C')]
A_hk, C_hk = grp(runs, 'Ahk'), grp(runs, 'Chk')
avg = lambda xs, k: sum(k(r) for r in xs) / len(xs)
cost = {}
for line in open(ex('real-project', 'cost.md'), encoding='utf-8'):
    m = re.match(r'\| (R[\w-]+) \| ([\d,]+) \| (\d+) \|', line)
    if m:
        cost[m.group(1)] = (int(m.group(2).replace(',', '')), int(m.group(3)))
ca = [v for k, v in cost.items() if k.split('-')[1].startswith('A')]
cc = [v for k, v in cost.items() if k.split('-')[1].startswith('C')]
k_up = sum(x[0] for x in cc) / len(cc) / (sum(x[0] for x in ca) / len(ca)) - 1
t_up = sum(x[1] for x in cc) / len(cc) / (sum(x[1] for x in ca) / len(ca)) - 1

# ---- demo 仓库新页面实验 ----
demo = json.load(open(ex('summary.json')))
dA = [r for r in demo if r['group'].startswith('A')]
dC = [r for r in demo if r['group'].startswith('C')]

# ---- 三个大型系统（figmafe 的输出，见 docs/real-test/large/）----
large = os.path.join(root, 'docs', 'real-test', 'large')
VERDICT = {'reuse': ('ok', '直接复用'), 'confirm': ('warn', '需确认'), 'new': ('bad', '需新建')}
SYSTEMS = [
    dict(id='yudao-mall', pf='h5', name='芋道商城', title='芋道商城 yudao-mall-uniapp', tech='H5 · uni-app', commit='4924ad51',
         stack='uni-app 商城前端：自有设计系统 <code>sheep/ui</code>（su-*）、业务组件 <code>sheep/components</code>（s-*），另装 uni-ui',
         first='三份设计稿<b class="bad">全部崩溃</b>；精修稿 12 ✓ / 1 △',
         fixes=['Vue props 里 <code>a ? b : c</code> 形式的默认值让 props 解析越界 → 冒号前必须有键名',
                'su-image 被当成 <code>&lt;image&gt;</code> 的封装，但项目直接写 <code>&lt;image&gt;</code> 71 次、su-image 4 次 → 封装组件用得不到原生一半时不算封装',
                '同一颜色有多个 token 时随手挑了一个 → 选项目里引用最多的 token',
                '卡片推荐了 SLiveCard（直播卡片）、资产卡片推荐了 SWalletCard → live / coupon / groupon / wallet 等商城领域词算专门用途；分割线别名加 subline（项目的分割线是 su-subline）',
                '体检把装修组件 s-notice-block（只是包了一层 su-notice-bar）算成第二种实现 → 记录组件模板里用到的组件，建立在同概念组件之上的算特化'],
         left='商品页底部操作栏 <code>detail-tabbar.vue</code>（4 个商品页在用）的按钮从插槽传入，文件本身看不出是操作栏，工具给的是 uni-ui GoodsNav（需确认）。'),
    dict(id='gkd', pf='android', name='gkd', title='gkd', tech='Android · Compose Multiplatform', commit='b2b42f28',
         stack='Jetpack Compose + Material3，自有 <code>ui/component</code>（GkTopAppBar、GkSettingItem、GkSwitch …）',
         first='三份设计稿 <b class="bad">5 ✗</b>；精修稿 <b class="bad">1 ✗</b>',
         fixes=['根目录只装了 husky 的 <code>package.json</code>，自动识别平台时被当成 H5 项目 → 依赖里有前端框架才算 H5',
                '按钮推荐了 GkIconButton（图标按钮 71 处），文字按钮其实是 TextButton（72 处）→ icon 算专门用途词',
                '搜索框是顶栏里的 GkAppBarTextField，名字里没有 search → 调用处大多绑定搜索词（<code>value = searchStr</code>）的输入组件按用法认作搜索框',
                '主图轮播推荐了 GkMatchingBanner（「开启匹配」的提示卡片）→ 没有图片 / 列表 / 页数参数的 banner 不算轮播',
                '体检「button 有 11 种实现」「dialog 有 15 种实现」「GkTwoLineText 是分割线」；RadioButton 被判成被 GkMenuItemRadioButton 封装 → 专门用途 / 弱词组件不算重复，按函数体记录「建立在谁之上」，封装组件 = 库组件名前加前缀和一个词'],
         left='涨跌标签推荐了 M3 InputChip（需确认）：项目里没有标签组件，严格说应新建。'),
    dict(id='androidproject', pf='android', name='AndroidProject', title='AndroidProject', tech='Android · XML View', commit='cdf35d36',
         stack='XML 骨架（国内很多 App 用它起步）：TitleBar 三方库 + 自有 <code>customWidget</code> 库（SettingBar、SwitchButton、SimpleRatingBar …）',
         first='三份设计稿 <b class="bad">8 ✗</b>；精修稿 <b class="bad">2 ✗</b>',
         fixes=['布局里 54 个 <code>AppCompatButton</code> 没算成平台 Button，推荐了只用过 1 次的 SubmitButton → AppCompat* 计入对应的平台控件',
                '设置行 SettingBar、状态布局 StatusLayout、评分 SimpleRatingBar 没被认出 → 补概念词，新增 rate / avatar 概念',
                'CountdownView 的注释写着「验证码倒计时」，是发送验证码的按钮 → 读组件声明前的文档注释，写明专门用途的降权，结论最多「需确认」',
                '头像推荐了 ScaleImageView（长按缩放的图片）→ 头像是独立概念，没有就判需新建'],
         left='支付倒计时其实要新建：工具给「需确认 CountdownView（专门用途：验证码）」，由人来否掉。'),
]
for sy in SYSTEMS:
    sy['reuse'] = json.load(open(os.path.join(large, sy['id'], 'reuse', sy['pf'] + '.json')))
    sy['sc'] = json.load(open(os.path.join(large, sy['id'], 'showcase', sy['pf'] + '.json')))
    md = open(os.path.join(large, sy['id'], 'reuse', sy['pf'] + '.md'), encoding='utf-8').read()
    sy['todo'] = re.findall(r'^\d+\. (.+)$', md.split('## 0.')[1].split('\n## ')[0], re.M) if '## 0.' in md else []
    w = re.search(r'\*\*设计系统封装\*\*[^：]*：(.+)', md)
    sy['wraps'] = w.group(1).replace('`', '') if w else ''
    sy['ev'] = [ev.get(sy['id'] + '/' + k) for k, _ in pages]
    sy['sum'] = {s: sum(v[s] for v in sy['ev'] if v) for s in '✓△✗'}


def md_code(t):
    return re.sub(r'`([^`]+)`', lambda m: '<code>' + m.group(1).replace('<', '&lt;') + '</code>', t)


def pick(it):
    return it['candidates'][0] if it['candidates'] and it['verdict'] != 'new' else None


def reuse_table(m):
    rows = []
    for it in m['items']:
        cls, label = VERDICT[it['verdict']]
        c = pick(it)
        why = (c['why'] if c else '项目和组件库里都没有').split('，⚠')[0]
        rows.append(f"<tr><td>{it['figma'].replace('（未组件化）', '')}</td><td class='{cls}'>{label}</td><td>{'<code>' + c['name'] + '</code>' if c else '—'}</td>"
                    f"<td>{why[:48] + ('…' if len(why) > 48 else '')}</td></tr>")
    return '<table class="sm"><tr><th style="width:17%">设计稿元素</th><th style="width:12%">结论</th><th style="width:26%">推荐</th><th>依据</th></tr>' + ''.join(rows) + '</table>'


def compare_table():
    head = '<tr><th style="width:16%">设计稿元素</th>' + ''.join(f"<th>{sy['name']}</th>" for sy in SYSTEMS) + '</tr>'
    rows = []
    for it in SYSTEMS[0]['sc']['items']:
        tds = []
        for sy in SYSTEMS:
            x = next(y for y in sy['sc']['items'] if y['figma'] == it['figma'])
            c = pick(x)
            tds.append(f"<td class='{VERDICT[x['verdict']][0]}'>{c['name'] if c else '需新建'}</td>")
        rows.append(f"<tr><td>{it['figma']}</td>{''.join(tds)}</tr>")
    return '<table class="sm">' + head + ''.join(rows) + '</table>'


def kpis(sy):
    r, me = sy['reuse'], sy['reuse']['metrics']
    return (f"<div class='kpi kpis'><div><span class='big'>{me['files']}</span> 个文件</div><div><span class='big'>{len(r['assets'])}</span> 个组件资产</div>"
            f"<div><span class='big'>{me['componentization']*100:.0f}%</span> 组件化率</div><div><span class='big'>{me['duplicateConcepts']}</span> 个重复概念</div></div>")


def system_block(sy, full):
    s = sy['sum']
    libs = '、'.join(f"{l['name']} {l['usage']} 次" for l in sy['reuse']['libraries'] if l['usage'] > 0)
    out = f"""<h3>{sy['title']}（{sy['tech']}） <span class="note">commit {sy['commit']} · 四份设计稿 ✓{s['✓']} △{s['△']} ✗{s['✗']}</span></h3>
<p class="note">{sy['stack']}</p>{kpis(sy)}
<p class="note">组件库用量：{libs}{('；设计系统封装：' + sy['wraps']) if sy['wraps'] else ''}</p>
<p class="lead">体检优先处理</p><ol class="note">{''.join('<li>' + md_code(t) + '</li>' for t in sy['todo'][:3])}</ol>"""
    if full:
        out += f"<p class='lead'>精修稿复用清单（直接复用 {sy['sc']['reuse']} · 需确认 {sy['sc']['confirm']} · 需新建 {sy['sc']['new']}）</p>" + reuse_table(sy['sc'])
    out += f"""<p class="lead">第一次跑（修正前）：{sy['first']}。出错的地方 → 通用修正</p>
<ul class="note">{''.join('<li>' + f + '</li>' for f in sy['fixes'])}</ul>
<p class="note"><b>还剩的问题</b>：{sy['left']}</p>"""
    return out


phones = ''.join(f"""<figure><img src="{rimg}showcase-{sy['id']}-frame.png"><figcaption>{sy['name']}（{sy['tech']}）<br>直接复用 {sy['sc']['reuse']} · 需确认 {sy['sc']['confirm']} · 需新建 {sy['sc']['new']}</figcaption></figure>""" for sy in SYSTEMS)
ev_head = '<tr><th>项目</th>' + ''.join(f'<th>{n}</th>' for _, n in pages) + '</tr>'
design_img = '../design/showcase/figma-export/product-detail/frame@2x.png'

page = f"""<!doctype html><html lang="zh"><head><meta charset="utf-8"><title>Figma → 代码：组件复用验证报告</title>
<style>
@page {{ size: A4; margin: 12mm 13mm; }}
* {{ box-sizing: border-box; }}
body {{ margin: 0; font-family: 'Noto Sans SC', sans-serif; color: #1f2329; font-size: 9.6pt; line-height: 1.55; }}
h1 {{ font-size: 20pt; margin: 0 0 1.5mm; letter-spacing: .5px; }}
h2 {{ font-size: 13pt; margin: 0 0 2.5mm; padding-left: 2.5mm; border-left: 3.5px solid #ff4d1a; }}
h3 {{ font-size: 10.8pt; margin: 3mm 0 1mm; }}
p {{ margin: 1.2mm 0; }}
.page {{ break-before: page; }} .page:first-child {{ break-before: auto; }}
.eyebrow {{ color: #ff4d1a; font-weight: 700; font-size: 8.6pt; letter-spacing: 1px; margin-bottom: 1mm; }}
.sub {{ color: #646a73; font-size: 9pt; margin-bottom: 3mm; }}
.lead {{ font-weight: 600; margin: 2mm 0 .8mm; }}
.box {{ background: #fff6f1; border: 1px solid #ffd8c8; border-radius: 6px; padding: 2.5mm 4mm; margin: 2mm 0; }}
.grey {{ background: #f5f6f7; border-radius: 6px; padding: 2.2mm 3.5mm; margin: 2mm 0; }}
table {{ border-collapse: collapse; width: 100%; font-size: 8.4pt; margin: 1.2mm 0; }}
th, td {{ border: 1px solid #e5e6eb; padding: .8mm 1.6mm; text-align: left; vertical-align: top; }}
th {{ background: #f5f6f7; font-weight: 600; }}
table.sm {{ font-size: 7.9pt; }} table.sm td, table.sm th {{ padding: .5mm 1.4mm; }}
.ev td, .ev th {{ white-space: nowrap; padding: .35mm 1.4mm; font-size: 8pt; }}
.ok {{ color: #00a86b; font-weight: 600; }} .bad {{ color: #e5383b; font-weight: 600; }} .warn {{ color: #d97706; }} .dim {{ color: #c0c4cc; }}
.tagho {{ font-size: 6.8pt; color: #fff; background: #8f959e; border-radius: 3px; padding: 0 1mm; }}
.flow {{ display: flex; gap: 1.5mm; margin: 2mm 0; }}
.flow div {{ flex: 1; background: #f5f6f7; border-radius: 5px; padding: 1.6mm 2mm; font-size: 8.2pt; }}
.flow b {{ display: block; color: #ff4d1a; font-size: 8.8pt; }}
.two {{ display: flex; gap: 5mm; align-items: flex-start; }} .two > * {{ flex: 1; min-width: 0; }}
.phones {{ display: flex; gap: 3mm; align-items: flex-start; margin: 2mm 0; }} .phones figure {{ flex: 1; margin: 0; min-width: 0; }}
.phones img {{ width: 100%; border: 1px solid #e5e6eb; border-radius: 6px; }}
figcaption {{ font-size: 7.8pt; color: #646a73; text-align: center; margin-top: .8mm; line-height: 1.35; }}
img.wide {{ width: 100%; border: 1px solid #e5e6eb; border-radius: 5px; }}
code {{ font-family: 'DejaVu Sans Mono', monospace; font-size: 7.8pt; background: #f3f4f6; padding: 0 .8mm; border-radius: 2px; }}
ul, ol {{ margin: 1mm 0; padding-left: 5mm; }} li {{ margin: .6mm 0; }}
.note {{ color: #4e5561; font-size: 8.4pt; }}
.big {{ font-size: 16pt; font-weight: 700; color: #ff4d1a; }}
.kpi {{ display: flex; gap: 2mm; margin: 2mm 0; }} .kpi div {{ flex: 1; border: 1px solid #ffd8c8; border-radius: 6px; padding: 1.8mm 2.5mm; font-size: 8.4pt; line-height: 1.4; }}
.kpis .big {{ font-size: 12.5pt; }} .kpi1 .big {{ font-size: 14pt; white-space: nowrap; }} .kpis div {{ padding: 1mm 2mm; }}
.foot {{ color: #8f959e; font-size: 7.8pt; margin-top: 3mm; }}
</style></head><body>

<section class="page">
<div class="eyebrow">验证报告 · 痛点 4「Figma 到实现存在偏差，组件未被充分复用」</div>
<h1>Figma → 代码：存量项目的组件复用</h1>
<div class="sub">仓库 luckylucky0001/figma-fe · 2026-10-20 · 交付物：CLI <code>figmafe</code>（Go）、Skill <code>figma-to-code</code>、免费的 Figma 导出插件、demo、分步教程、演示 deck、本报告</div>
<div class="two">
<div style="flex:2.5">
<div class="kpi kpi1">
<div><span class="big">{tot['✓']}/{n_items}</span><br>复用推荐正确，<b>{tot['✗']}</b> 处错误</div>
<div><span class="big">{len(projects)}</span><br>个真实开源项目，H5 / Android 主流技术栈</div>
<div><span class="big">{passes(A_all)} → {passes(C_all)}</span><br>Agent 实现新页面通过复用检查（不用 → 用 Skill）</div>
</div>
<div class="box"><b>结论</b><ul>
<li><b>先把「项目里有什么」查清楚，才知道新设计稿能复用什么。</b><code>figmafe reuse</code> 扫一遍现有项目，得到组件资产、项目惯例、重复实现，以及「组件库里有、但本项目用不了」的组件；<code>figmafe match</code> 再把设计稿逐个元素对上去，给出直接复用 / 需确认 / 需新建，附项目里的真实写法和 props 配置。</li>
<li><b>推荐准不准</b>：{len(projects)} 个项目、4 份设计稿共 {n_items} 项，✓{tot['✓']} / △{tot['△']} / ✗{tot['✗']}。其中 {len(heldout)} 个项目是留出的：调规则时没用过，第一次跑出过错（含三个大型系统：芋道商城、gkd、AndroidProject），修的都是通用规则（第 4、5 节）。</li>
<li><b>Agent 用不用 Skill 差别很大</b>：5 个真实项目里按同一份设计稿实现新页面，不用 Skill {passes(A_all)} 通过复用检查，用 Skill {passes(C_all)}；代价是多 {k_up*100:.0f}% token、{t_up*100:.0f}% 耗时（第 8 节）。</li>
<li><b>真的跑得起来</b>：在 vue-h5-template 里 Agent 按 Skill 做出新页面，复用门禁、类型检查、构建、单测都通过；像素对比因为项目主题和设计稿不同没过，这正是存量项目要处理的事（第 8 节）。</li>
<li><b>只用免费能力</b>：不用 Code Connect、Dev Mode MCP、Variables REST API；设计数据由自己写的插件（Plugin API）导出。</li>
</ul></div>
<p class="note">右图是本报告的演示设计稿（第 3 节）：一张完成度较高的商品详情页，含主图轮播、秒杀价格、规格配送、评价、底部购买栏。没有真实 Figma 文件可用，它和另外三份设计稿一样由手写的设计源导出成插件格式（局限见第 9 节）。</p>
</div>
<div><img class="wide" src="{design_img}"></div>
</div>
<h3>交付物</h3>
<table>
<tr><th style="width:24%">交付物</th><th style="width:30%">位置</th><th>说明</th></tr>
<tr><td>CLI <code>figmafe</code>（Go）</td><td><code>cli/</code>，构建到 <code>bin/figmafe</code></td><td><code>reuse</code> 体检、<code>match</code> 复用清单、<code>compare</code> / <code>verify</code> 截图与门禁，单二进制，无外部依赖</td></tr>
<tr><td>Skill <code>figma-to-code</code></td><td><code>.claude/skills/figma-to-code/</code></td><td>固定 Agent 的流程：体检 → 复用清单 → 人确认 → 实现 → 检查到通过</td></tr>
<tr><td>Figma 导出插件</td><td><code>figma-plugin/</code></td><td>只用免费的 Plugin API，导出节点树、组件集、变体和 @2x 图</td></tr>
<tr><td>Demo</td><td><code>h5/</code>、<code>android/</code>、<code>design/</code></td><td>H5（React）与 Android（Compose / XML）组件库、四份设计稿（含精修的商品详情页）</td></tr>
<tr><td>教程、报告、deck</td><td><code>docs/tutorial/</code>、<code>docs/report.pdf</code></td><td>分步教程带截图；本报告；5 页演示 deck</td></tr>
<tr><td>实验数据</td><td><code>experiments/</code>、<code>docs/real-test/</code></td><td>26 个项目的回归评测、Agent A / C 实验、三个大型系统的完整输出</td></tr>
</table>
</section>

<section class="page">
<h2>1. 问题和方案</h2>
<p>新页面从 Figma 落到代码时，偏差主要不在像素，而在<b>没用上项目已有的组件</b>：Agent 或新人看不到项目里有什么，就照着设计稿重写一套导航栏、卡片、按钮，再写死颜色。项目越大越严重，而大项目恰恰最需要复用。所以方案的核心是先盘点存量，再对设计稿：</p>
<div class="flow">
<div><b>① 体检存量</b><code>reuse</code>：组件资产、项目惯例、重复概念、复制粘贴结构、写死颜色、版本 / 主题兼容性</div>
<div><b>② 复用清单</b><code>match</code>：设计稿元素 → 项目组件，附真实用法、props 配置、可视化叠图</div>
<div><b>③ 人确认</b>「需确认 / 需新建」写进 <code>component-map.yaml</code>，下次自动采用</div>
<div><b>④ 写完检查</b><code>reuse --focus / --changed</code>：自造轮子、有组件没用、手写已有组件</div>
<div><b>⑤ 对照设计</b><code>compare</code> / <code>verify</code>：截图、像素与节点级偏差，带 Figma 节点 id</div>
<div><b>⑥ 不许变差</b><code>--baseline --no-regress</code> 接 CI；经验沉淀进 <code>rules.md</code></div>
</div>
<table>
<tr><th style="width:23%">Proposal 要求</th><th>落地方式</th></tr>
<tr><td>Figma 组件与代码组件的对应关系</td><td><code>match</code> 按名字、UI 概念（中英文同义词）、项目前缀、使用次数、设计系统封装、主用组件库等打分；人确认后写进 <code>component-map.yaml</code>，之后优先采用确认过的映射</td></tr>
<tr><td>开发前识别可复用组件</td><td>先 <code>reuse</code> 再 <code>match</code>：只推荐<b>本项目真能用</b>的组件（检查组件库版本、改名 / 移除、主题依赖、废弃 API、是否被设计系统封装）；没有组件但有固定写法时推荐「项目惯例」</td></tr>
<tr><td>开发后检查用了哪些组件、和 Figma 截图对比</td><td><code>reuse --focus &lt;页面&gt;</code> 或 <code>--changed &lt;git ref&gt;</code> 只看新增代码；<code>compare</code> 截图对比；demo 里 <code>verify</code> 把两者合成一个门禁</td></tr>
<tr><td>偏差反馈给 Agent、沉淀规则</td><td>每条问题带文件行号或 Figma 节点 id 和修改建议，Skill 规定照着改到通过；决定和理由写进 component-map，经验写进 <code>rules.md</code></td></tr>
</table>
<h3>Skill 把流程固定下来</h3>
<p><code>figma-to-code</code> Skill 规定 Agent 的顺序：<b>先体检，再出复用清单并等人确认，再实现，最后跑页面检查和截图对比直到通过</b>，不允许跳过前两步直接写页面。第 8 节的实验说明这个顺序比模型本身更重要：同样的模型，按 Skill 走和不按 Skill 走，复用检查通过率是 {passes(C_all)} 对 {passes(A_all)}。</p>
<h3>支持的技术栈</h3>
<table>
<tr><th style="width:12%">平台</th><th style="width:30%">框架</th><th>组件库 / 设计系统</th></tr>
<tr><td>H5</td><td>Vue 2 / 3、React、uni-app、Taro</td><td>Vant 2 / 4、React Vant、NutUI、antd-mobile、uni-ui、uView、wot-design-uni、taro-ui、TDesign、Element Plus，以及项目自有组件</td></tr>
<tr><td>Android</td><td>XML View（含 DataBinding）、Jetpack Compose（含 Multiplatform）</td><td>Material Components / Material3、AppCompat、QMUI、三方 View 库（按包名识别），以及项目自有设计系统</td></tr>
</table>
<div class="grey note"><b>只用免费能力</b>：Figma 的 Code Connect、Dev Mode MCP 和 Variables REST API 都要付费，这里一个都不用。设计数据由自己写的 Figma 插件导出（Plugin API 免费），格式见 <code>docs/CONTRACTS.md</code>。</div>
</section>

<section class="page">
<h2>2. 工具做了什么</h2>
<h3>2.1 体检：组件库有，不等于本项目能用</h3>
<div class="two">
<div><p><code>figmafe reuse &lt;仓库&gt;</code> 只读扫描，输出：</p><ul class="note">
<li><b>组件资产</b>：组件库组件、设计系统、业务组件、自定义 View / Composable，每个都带使用次数和项目里的真实写法</li>
<li><b>主用组件库</b>：装了多套 H5 组件库时，按 <code>.env</code> 里的 UI 配置项或使用次数定一套，新页面不许混用</li>
<li><b>兼容性</b>：组件库版本里没有的组件、主题不支持会崩溃的组件、已废弃的 API、被设计系统封装过的库组件</li>
<li><b>项目惯例</b>：没有组件但固定这么写的结构（如 1px 细线 View）</li>
<li><b>重复概念、复制粘贴结构、写死颜色</b>（含已有同值 token 的），按影响面排出「优先处理」</li>
</ul></div>
<div><img class="wide" src="{img}21-qmui-compat.png"><p class="note">QMUI_Android 用 AppCompat 主题和 Material 1.4：列出的组件会崩溃或编译失败，match 不会推荐。</p></div>
</div>
<h3>2.2 复用清单：依据是项目自己的用法</h3>
<p><code>figmafe match --figma &lt;导出目录&gt; &lt;仓库&gt;</code> 对设计稿里每个组件实例或未组件化的图层打分，信号包括：</p>
<table class="sm">
<tr><th style="width:22%">信号</th><th>例子</th></tr>
<tr><td>名字与 UI 概念</td><td>NavBar / TitleBar / TopAppBar 同为导航栏；中文图层「持仓列表项」→ 列表项；ProductCard → Card 算匹配，MaterialRadioButton 不算 Button</td></tr>
<tr><td>项目用得多、设计系统优先</td><td>芋道的 SuNavbar（设计系统）优先于 uni-ui NavBar；gkd 的 TextButton（72 处）优先于 Button（3 处）</td></tr>
<tr><td>专门用途降权</td><td>名字或<b>文档注释</b>写明一种用途的组件（PasswordInputView、GkIconButton、注释为「验证码倒计时」的 CountdownView）不当通用组件</td></tr>
<tr><td>按用法认概念</td><td>调用处都在传搜索词的输入组件是搜索框；含 <code>&lt;router-view&gt;</code> 的组件是页面外壳</td></tr>
<tr><td>变体 → props</td><td>Tone=Warning → <code>type="warning"</code>；尺寸按设计稿高度选（32px 的 Medium 按钮 → Vant <code>size="small"</code>）</td></tr>
</table>
<h3>2.3 写完检查：有现成组件却手写</h3>
<div class="two"><div><img class="wide" src="{img}26-check-a.png"></div>
<div><p class="note">不用 Skill 的实现（newbee，A 组）：分割线、商品卡、底部操作栏都是手写的，每条带行号和应复用的组件。<code>--changed origin/main</code> 只检查新增代码，存量问题不挡提交；<code>--baseline --no-regress</code> 让组件化率、写死颜色等指标只许变好。</p></div></div>
</section>

<section class="page">
<h2>3. 演示：同一张设计稿，三个大型系统各自的复用清单</h2>
<p>同一张商品详情设计稿，放到三个技术栈完全不同的系统里，推荐的是<b>各自项目里的组件</b>：芋道用自己的 SuSwiper / SuNavbar 和 uni-ui Rate，gkd 用 GkTopAppBar / GkSettingItem 和 Compose 的 HorizontalPager，AndroidProject 用 TitleBar / SettingBar 和自己的 SimpleRatingBar。绿框直接复用、橙框需确认、红框需新建。</p>
<div class="phones">
<figure><img src="{design_img}"><figcaption>设计稿（商品详情）</figcaption></figure>
{phones}
</div>
{compare_table()}
<p class="note">设计稿源文件 <code>design/showcase/source/product-detail.html</code>，导出 <code>design/showcase/figma-export/product-detail/</code>；三个系统的完整输出在 <code>docs/real-test/large/&lt;系统&gt;/showcase/</code>。</p>
</section>

<section class="page">
<h2>4. 推荐准不准：{len(projects)} 个真实开源项目</h2>
<p><b>方法</b>：每个项目固定 commit，可复现；先逐条写好标准答案（✓ 正确、△ 可用但不是最优、✗ 错误，「—」表示应该新建），再跑工具对答案。前三份设计稿覆盖组件化 / 未组件化、中文 / 英文图层和新概念，在全部 {len(projects)} 个项目上判分；第四份精修稿（商品详情）在 {len(sc_projects)} 个项目上判分。<b>留出</b>项目是调规则时没用过的，第一次跑的结果见第 5 节。</p>
<div class="two"><table class="ev">{ev_head}
{eval_rows(0)}
</table><table class="ev">{ev_head}
{eval_rows(1)}
</table></div>
<div class="box"><b>合计 ✓{tot['✓']} / △{tot['△']} / ✗{tot['✗']}</b>（{n_items} 项）。△ 多数是工具给了「需确认」、由人决定的情况，例如项目里没有标签组件时推荐组件库的 Chip，或只有别的组件库里有的 Price。</div>
<p class="note">技术栈：Vue 2 / 3、React、uni-app、Taro × Vant / NutUI / antd-mobile / uni-ui / taro-ui / wot-design-uni / 自有组件；Android XML / Compose × QMUI / Material / AppCompat / 三方 View 库 / 自有设计系统。逐项结果见 <code>experiments/real-match-eval.md</code>，每一条的判断理由见 <code>match-judgement.md</code>。回归：<code>python3 experiments/real-match-eval.py</code>。</p>
</section>

<section class="page">
<h2>5. 边测边改：第一次跑的错误和通用修正</h2>
<p>每加一批项目，先按写好的答案跑一次、如实记录错误，再改成<b>通用规则</b>（没有针对单个项目的特判），最后对全部项目回归，并重算真实项目实验的打分（<code>experiments/rescore-real.sh</code>）。</p>
<table>
<tr><th style="width:22%">批次 / 项目</th><th style="width:13%">第一次跑</th><th>出错原因 → 通用修正</th></tr>
<tr><td>Jetsnack</td><td><b class="bad">3 ✗</b></td><td>桌面小部件（Glance）组件被当成 App 组件；没认出 <code>Jetsnack*</code> 项目前缀 → 识别 Glance、学习项目前缀、<code>ui/components</code> 算设计系统</td></tr>
<tr><td>litemall</td><td><b class="bad">2 ✗</b></td><td>推荐了 Vant 3 才有的 ActionBar，项目是 Vant 2 → 组件库加入改名 / 移除信息，按实际安装版本过滤</td></tr>
<tr><td>hello-uniapp</td><td><b class="bad">12 ✗</b></td><td>uni-app 内置小写标签被当成项目组件、中文图层名规范化后为空被判「同名」→ 加 uni-ui 目录，修空名误判</td></tr>
<tr><td>NewPipe</td><td><b class="bad">3 ✗</b></td><td>AppCompat 标准控件没算进项目资产，XML 页面被推荐 Compose 组件 → 跨技术栈的候选最多「需确认」</td></tr>
<tr><td>mall-app-web / AntennaPod</td><td>崩溃 / <b class="bad">1 ✗</b></td><td>uni-app 内置标签没有 import 路径导致崩溃；AntennaPod 的空状态封装 EmptyViewHandler 不是 View → 修崩溃、识别「界面封装类」</td></tr>
<tr><td>WanAndroid / Seal / KunMinX / jeecg</td><td><b class="bad">11 / 7 / 3 / 11 ✗</b></td><td>三方 View 库和 wot-design-uni 没被收录；Android 的 ActionBar 是顶栏；PasswordInputView / ConfirmButton 这类专门用途组件被当成通用组件；会崩溃的组件没有别的候选时仍排第一</td></tr>
<tr><td>买入下单（第三份设计稿，14 个项目）</td><td><b class="bad">12 ✗</b></td><td>开关 / 复选框 / 单选原来是一个概念，搜索算在输入框里 → 拆开概念、补词表</td></tr>
<tr><td>vue-h5-template</td><td><b class="bad">23 ✗</b></td><td>同时装了 Vant / NutUI / Varlet，打平后一半推荐 NutUI，同一页面混用 → 先定「主用组件库」（看 <code>.env</code>）；演示页被当成设计系统 → 目录按最外层分类</td></tr>
<tr><td>芋道 / gkd / AndroidProject</td><td>崩溃 / <b class="bad">5 / 8 ✗</b></td><td>props 三元默认值让扫描崩溃；图标按钮当文字按钮；54 个 AppCompatButton 漏数；「验证码倒计时」被当成倒计时（详见第 6、7 节）</td></tr>
<tr><td>商品详情（精修稿，4 个项目）</td><td><b class="bad">3 ✗</b></td><td>提示卡片 GkMatchingBanner 因为名字含 banner 被当成轮播；评分控件、头像没有概念 → 新增 rate / avatar 概念，没有图片参数的 banner 不算轮播</td></tr>
</table>
<h3>准确率随项目增加的变化</h3>
<table>
<tr><th style="width:24%">阶段</th><th style="width:12%">项目</th><th style="width:20%">结果</th><th>这一阶段新增的通用规则（举例）</th></tr>
<tr><td>报告初版</td><td>22</td><td>✓699 △27 ✗0</td><td>版本 / 主题兼容性、项目惯例、设计系统封装</td></tr>
<tr><td>+ vue-h5-template</td><td>23</td><td>✓734 △25 ✗0</td><td>主用组件库；只有依赖里的包才算三方库</td></tr>
<tr><td>+ 三个大型系统</td><td>26</td><td>✓826 △32 ✗0</td><td>读组件文档注释；AppCompat* 算平台控件；建立在其他组件之上的算特化</td></tr>
<tr><td>+ 实现反馈</td><td>26</td><td>✓827 △31 ✗0</td><td>按设计稿高度选尺寸；按调用处认出搜索框</td></tr>
<tr><td>+ 精修演示稿</td><td>26</td><td>✓{tot['✓']} △{tot['△']} ✗{tot['✗']}</td><td>评分、头像概念；没有内容参数的 banner 不是轮播</td></tr>
</table>
<p class="note">vant-demo、Pokedex、Sunflower、vue3-vant-mobile 第一次跑就没有 ✗。真实项目实验的打分在这几轮修改中都没有变化。</p>
</section>

<section class="page">
<h2>6. 大型系统：H5</h2>
<p>前面的项目多是模板和示例，所以又挑了三个规模大、写法杂的生产系统，按 Skill 的顺序分析：体检 → 复用清单 → 人工逐条判断 → 出错的改成通用规则 → 全量回归。答案在跑当前版本之前写好，「第一次跑」是修正前的版本。</p>
{system_block(SYSTEMS[0], True)}
</section>

<section class="page">
<h2>7. 大型系统：Android（Compose 与 XML）</h2>
{system_block(SYSTEMS[1], False)}
{system_block(SYSTEMS[2], False)}
</section>

<section class="page">
<h2>8. 实现：Agent 照着清单写页面</h2>
<h3>8.1 vue-h5-template：装依赖、类型检查、构建、像素对比都跑了</h3>
<div class="flow">
<div><b>① 体检</b>装了 Vant / NutUI / Varlet，<code>.env</code> 里 <code>VITE_UI_FRAMEWORK=vant</code> → 主用 Vant</div>
<div><b>② 复用清单</b>第一次 10 ✓ / 23 ✗（混用两套组件库）→ 修正后 32 ✓ / 1 △</div>
<div><b>③ 实现</b><code>src/views/order/detail.vue</code> + 新建 <code>PriceText.vue</code>，理由写进 component-map</div>
<div><b>④ 验收</b>复用门禁 · <code>vue-tsc</code> · <code>vite build</code> · 50 个单测，全部通过</div>
</div>
<p><b>像素对比 7.11%，没通过</b>：项目主题（靛蓝主色、44px 按钮、字体）和设计稿不同。组件复用对了，样式还要对齐项目的主题 token。Agent 反馈的工具问题都已修：尺寸按设计稿高度选（32px 的 Medium 按钮 → <code>size="small"</code>，而不是按名字对到 44px 的 normal）；含 <code>&lt;router-view&gt;</code> 的 Layout 不再当成操作栏；CellGroup 的 title 显示在卡片外这一点写进了提示。</p>
<h3>8.2 用不用 Skill：5 个真实项目的 A / C 对照</h3>
<p class="note">同一份设计稿（订单详情）。A 组只给任务和 Figma 导出，C 组加上 Skill。所有 run 用同一版 <code>figmafe reuse --focus / --changed</code> 打分：自造已有组件、有组件没用、手写已有组件，任何一项都算不通过。</p>
<table>
<tr><th rowspan="2">项目</th><th rowspan="2">技术栈</th><th colspan="3">A 不用 Skill（Opus / Haiku）</th><th colspan="3">C 用 Skill（Opus / Haiku）</th></tr>
<tr><th>通过</th><th>复用组件数</th><th>问题数</th><th>通过</th><th>复用组件数</th><th>问题数</th></tr>
{exp_rows()}
<tr><th colspan="2">合计</th><th class="bad">{passes(A_all)}</th><th>平均 {avg(A_all, lambda r: r['used']):.1f}</th><th>共 {sum(bad(r) for r in A_all)}</th><th class="ok">{passes(C_all)}</th><th>平均 {avg(C_all, lambda r: r['used']):.1f}</th><th>共 {sum(bad(r) for r in C_all)}</th></tr>
</table>
<ul class="note">
<li><b>Now in Android</b>：A 组没先查已有组件，新建了 NiaNavBar / NiaCard / NiaDivider / NiaCell 等，和已有的 NiaTopAppBar / Card / HorizontalDivider / ListItem 重复，还写了一套 18 个写死颜色的主题。C 组复用了它们，新建的组件都写了理由。</li>
<li><b>Haiku</b>：从 {passes(A_hk)} 提到 {passes(C_hk)}，模型越弱，固定流程的收益越大。</li>
<li><b>成本</b>：有用量记录的 {len(ca)} 次 A、{len(cc)} 次 C，C 平均多 {k_up*100:.0f}% token、{t_up*100:.0f}% 时间，多花在体检、复用清单、页面检查和写决定上。</li>
</ul>
</section>

<section class="page">
<h2>9. 新页面闭环、局限和接入</h2>
<h3>9.1 demo 仓库：一条命令做门禁</h3>
<p>自建的 H5（React）和 Android（Compose，兼容 compileSdk 32 / minSdk 23，另有 XML View 存量组件）组件库加一份订单详情设计稿，跑了 {len(demo)} 次 A / C 对照。<code>figmafe verify</code> 一条命令输出复用率、违规、像素差异和节点级偏差，退出码可做 CI 门禁，能区分「页面写错」和「组件本身不符合设计」（<code>docs/component-issues.md</code>）；教程第 8 步故意写错 3 处，全部被抓出。这个仓库小、组件好找，A 组复用率也是 100%，<b>证明不了复用率提升</b>；差别在规范：A 组 {sum(r['selfBuiltCountdown'] for r in dA)}/{len(dA)} 次自造了设计稿外的 CountDown，C 组 {sum(r['selfBuiltCountdown'] for r in dC)}/{len(dC)}。</p>
<div class="two">
<div><img class="wide" src="{img}09-verify-fail.png"><p class="note">教程里故意写错的版本：每条问题带 Figma 节点 id 和修改建议。</p></div>
<div><img class="wide" src="{img}10-verify-pass.png"><p class="note">照着改完，<code>verify</code> 通过。</p></div>
</div>
<h3>9.2 局限（如实说明）</h3>
<ul class="note">
<li><b>没有真实 Figma 文件</b>：本环境访问不了 Figma，四份设计稿都由手写的设计源导出成插件格式；插件只用模拟的 Plugin API 测过，没在真实 Figma 里跑过。</li>
<li><b>大多没有编译</b>：除了 vue-h5-template，Agent 实验只做静态检查，没有构建和运行这些开源项目（本环境装不了 Android SDK）；demo 的 Android 截图用 Compose Desktop 渲染，不是真机。</li>
<li><b>实验有泄漏</b>：A 组 agent 也能看到本仓库的 Skill，1 次 A run 自己跑了 figmafe，已剔除；泄漏只会让 A 组更好，所以 A / C 差距是保守估计。样本量小（每格 1～2 次），只能说明方向。</li>
<li><b>评分标准是人定的</b>：✓ / △ / ✗ 由我逐条判断并写成回归用例，真实接入时应由组件 owner 复核一次。</li>
<li><b>只看静态代码</b>：用途只能从插槽或运行时看出的组件（芋道的 detail-tabbar）认不出；组件库尺寸表目前只收录了 Vant / NutUI 的按钮。</li>
</ul>
<h3>9.3 接入你们的项目</h3>
<table>
<tr><th style="width:4%">#</th><th>动作</th><th style="width:36%">产出</th></tr>
<tr><td>1</td><td><code>figmafe reuse &lt;H5 / Android 仓库&gt;</code>（只读，不改代码）</td><td>资产清单、优先处理项、首份体检基线</td></tr>
<tr><td>2</td><td>Figma 插件导出下一个要做的页面，<code>figmafe match</code></td><td>复用清单 + 叠图，组件 owner 确认「需确认」项</td></tr>
<tr><td>3</td><td>开发时让 Agent 按 Skill 走；提交前 <code>reuse --changed origin/main</code></td><td>新增代码的复用问题，带行号</td></tr>
<tr><td>4</td><td>CI 里 <code>reuse --baseline main.json --no-regress</code>（示例：<code>docs/ci/</code>）</td><td>组件化率、写死颜色等指标只许变好</td></tr>
</table>
<p class="foot">教程 <code>docs/tutorial/README.md</code> · 数据 <code>experiments/real-match-eval.md</code>、<code>match-judgement.md</code>、<code>real-project/summary.md</code> · 重新生成 <code>python3 docs/tools/build_report.py</code></p>
</section>
</body></html>"""
open(os.path.join(root, 'docs', 'report.html'), 'w').write(page)
print('✓ docs/report.html')
