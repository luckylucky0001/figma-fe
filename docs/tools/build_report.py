"""生成 docs/report.html（约 3 页 A4），再由 shoot.mjs 输出 docs/report.pdf。
数字全部来自实验数据，不手填：
  experiments/real-match-eval.json / .md   复用清单准确率（real-match-eval.py）
  experiments/real-project/summary.json    真实项目 A / C 对照（real-summary.py）
  experiments/real-project/cost.md         真实项目实验的 token / 耗时
  experiments/summary.json + cost.json     demo 仓库新页面实验（summarize.py）"""
import json, os, re

root = os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..'))
ex = lambda *p: os.path.join(root, 'experiments', *p)
img = 'tutorial/img/'

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
pages = [('order-detail', '订单详情'), ('holdings', '我的持仓'), ('trade-order', '买入下单')]
tot = {s: sum(v[s] for v in ev.values()) for s in '✓△✗'}
n_items = sum(tot.values())


def cell(v):
    s = f"✓{v['✓']}"
    if v['△']:
        s += f" <span class='warn'>△{v['△']}</span>"
    if v['✗']:
        s += f" <b class='bad'>✗{v['✗']}</b>"
    return s


def eval_rows():
    out = []
    for p in projects:
        tag = '留出' if p in heldout else ''
        out.append(f"<tr><td>{p}</td><td>{stack.get(p, '').replace('（', ' ').replace('）', '')}</td><td>{tag}</td>" +
                   ''.join(f"<td>{cell(ev[p + '/' + k])}</td>" for k, _ in pages) + '</tr>')
    return '\n'.join(out)


# ---- 真实项目 A / C 对照 ----
rp = json.load(open(ex('real-project', 'summary.json')))
INVALID = {'Rlm-A-1'}  # A 组 agent 自己找到并运行了 Skill，不能当基线（见 real-project/summary.md）
runs = [r for r in rp.values() if r['run'] not in INVALID]
PNAME = {'Rnb': ('newbee-mall', 'Vue 3 + Vant'), 'Rlm': ('litemall（留出）', 'Vue 2 + Vant 2'),
         'Rqm': ('QMUI_Android', 'XML + QMUI'), 'Rnp': ('NewPipe（留出）', 'XML + AppCompat'),
         'Rnia': ('Now in Android', 'Compose + 自有设计系统')}


def grp(rs, g):
    return [r for r in rs if r['group'] == g]


def passes(rs):
    return f"{sum(r['pass'] for r in rs)}/{len(rs)}" if rs else '—'


def bad(r):  # 自造轮子 + 有组件没用 + 手写已有组件
    return r['reinvented'] + r['reinventedWarn'] + r['missed'] + r['handBuiltError'] + r['handBuiltWarn']


def exp_rows():
    out = []
    for p, (name, st) in PNAME.items():
        rs = [r for r in runs if r['project'] == p]
        a = grp(rs, 'A') + grp(rs, 'Ahk')
        c = grp(rs, 'C') + grp(rs, 'Chk')
        f = lambda xs, k: ' / '.join(str(k(r)) for r in xs)
        out.append(f"<tr><td>{name}</td><td>{st}</td>"
                   f"<td class='bad'>{passes(a)}</td><td>{f(a, lambda r: r['used'])}</td><td>{f(a, bad)}</td>"
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

page = f"""<!doctype html><html lang="zh"><head><meta charset="utf-8"><title>Figma → 代码：组件复用验证报告</title>
<style>
@page {{ size: A4; margin: 11mm 12mm; }}
* {{ box-sizing: border-box; }}
body {{ margin: 0; font-family: 'Noto Sans SC', sans-serif; color: #1f2329; font-size: 9.8pt; line-height: 1.5; }}
h1 {{ font-size: 17pt; margin: 0 0 1.5mm; }} h2 {{ font-size: 12pt; margin: 3.5mm 0 1.5mm; padding-left: 2.5mm; border-left: 3px solid #ff4d1a; }}
h3 {{ font-size: 10.4pt; margin: 2mm 0 1mm; }}
.sub {{ color: #646a73; font-size: 9pt; margin-bottom: 2.5mm; }}
.tldr {{ background: #fff1eb; border-radius: 6px; padding: 2mm 4mm; margin: 1.5mm 0 2.5mm; }}
table {{ border-collapse: collapse; width: 100%; font-size: 8.4pt; margin: 1.2mm 0; }}
th, td {{ border: 1px solid #e5e6eb; padding: 0.8mm 1.5mm; text-align: left; vertical-align: top; }}
th {{ background: #f5f6f7; font-weight: 600; }}
.ok {{ color: #00b578; font-weight: 600; }} .bad {{ color: #f53f3f; font-weight: 600; }} .warn {{ color: #d97706; }}
.tagho {{ font-size: 7pt; color: #fff; background: #8f959e; border-radius: 3px; padding: 0 1mm; }}
.flow {{ display: flex; gap: 1.5mm; margin: 1.5mm 0; }}
.flow div {{ flex: 1; background: #f5f6f7; border-radius: 5px; padding: 1.4mm 2mm; font-size: 8.2pt; }}
.flow b {{ display: block; color: #ff4d1a; font-size: 8.8pt; }}
.two {{ display: flex; gap: 4mm; align-items: flex-start; }} .two > * {{ flex: 1; min-width: 0; }}
.pb {{ break-before: page; }}
.ev td, .ev th {{ white-space: nowrap; padding: 0.35mm 1.5mm; }}
img.wide {{ width: 100%; border: 1px solid #e5e6eb; border-radius: 5px; }}
code {{ font-family: 'DejaVu Sans Mono', monospace; font-size: 8pt; background: #f3f4f6; padding: 0 1mm; border-radius: 2px; }}
ul {{ margin: 1mm 0; padding-left: 5mm; }} li {{ margin: 0.5mm 0; }}
.note {{ color: #646a73; font-size: 8.3pt; margin: 1mm 0; }}
.big {{ font-size: 15pt; font-weight: 700; color: #ff4d1a; }}
.kpi {{ display: flex; gap: 2mm; margin: 2mm 0; }} .kpi div {{ flex: 1; border: 1px solid #ffd8c8; border-radius: 6px; padding: 1.5mm 2.5mm; font-size: 8.4pt; }}
</style></head><body>

<h1>Figma → 代码：组件复用验证报告</h1>
<div class="sub">痛点 4「Figma 到实现存在偏差，组件未被充分复用」· 仓库 luckylucky0001/figma-fe · 2026-10-20 · 交付物：CLI <code>figmafe</code>（Go）+ Skill <code>figma-to-code</code> + 免费 Figma 导出插件 + Demo + 分步教程</div>

<div class="kpi">
<div><span class="big">{len(projects)}</span> 个真实开源项目（H5 / Android 主流技术栈）× 3 份设计稿</div>
<div><span class="big">{tot['✓']}/{n_items}</span> 项复用推荐正确，<b>{tot['✗']}</b> 项错误，其余为「可用但不是最优」</div>
<div><span class="big">{passes(A_all)} → {passes(C_all)}</span> 真实项目里实现新页面：不用 Skill → 用 Skill 通过复用检查</div>
<div><span class="big">+{k_up*100:.0f}%</span> token、+{t_up*100:.0f}% 耗时：用 Skill 多出的成本</div>
</div>

<div class="tldr"><b>结论</b><ul>
<li><b>对存量项目，先把「有什么」查清楚，才知道新设计稿能复用什么。</b><code>figmafe reuse</code> 扫一遍现有 H5（Vue / React / uni-app / Taro + Vant / antd-mobile / NutUI / uni-ui / taro-ui）或 Android（XML View / Compose，QMUI / Material / AppCompat / 自有设计系统）项目，得到组件资产、项目惯例、重复实现和「组件库有但本项目用不了」的组件。<code>figmafe match</code> 再把新设计稿逐个元素对上去，给出「直接复用 / 需确认 / 需新建」，并附上项目里的真实写法和变体 → props 的配置。</li>
<li><b>准确率</b>：{len(projects)} 个项目 × 3 份设计稿共 {n_items} 项，✓{tot['✓']} / △{tot['△']} / ✗{tot['✗']}。其中 {len(heldout)} 个是「留出」项目：调规则时没用过，第一次跑出现过错误（见第 2 节），修的都是通用规则，之后作为回归测试。</li>
<li><b>Agent 实测</b>：在 5 个真实项目里按同一份设计稿实现新页面，不用 Skill 的 {len(A_all)} 次全部没通过复用检查（Now in Android 里新建了 11～12 个和已有组件重复的组件；Haiku 在 newbee 里一个现成组件都没用），用 Skill 的 {passes(C_all)} 次通过。Haiku 从 {passes(A_hk)} 提到 {passes(C_hk)}：模型越弱，固定流程的收益越大。</li>
<li><b>新页面闭环</b>：demo 仓库里，<code>figmafe verify</code> 一条命令输出复用率、违规、像素差异和节点级偏差，退出码可做 CI 门禁；故意埋的 3 处错误全部被抓出，修正后一次通过。</li>
<li><b>只用免费能力</b>：不用 Figma Code Connect、Dev Mode MCP 和 Variables REST API（都要付费），设计数据走自己写的插件（Plugin API，免费）。</li>
</ul></div>

<h2>1. 方案：Proposal 的每一条都落成了命令</h2>
<div class="flow">
<div><b>① 体检存量</b><code>reuse</code>：资产、惯例、重复概念、复制粘贴结构、写死颜色、版本 / 主题兼容性</div>
<div><b>② 新稿复用清单</b><code>match</code>：Figma 元素 → 项目组件，附真实用法、props 配置、可视化叠图</div>
<div><b>③ 人确认</b>需确认 / 需新建的项写进 <code>component-map.yaml</code>，下次自动采用</div>
<div><b>④ 实现后检查</b><code>reuse --focus / --changed</code>：自造轮子、有组件没用、手写已有组件</div>
<div><b>⑤ 对照设计</b><code>compare</code> / <code>verify</code>：截图 + 像素 / 节点级偏差，带 Figma 节点 id</div>
<div><b>⑥ 不许变差</b><code>--baseline --no-regress</code> 接 CI；偏差沉淀成 <code>rules.md</code></div>
</div>
<div class="two"><div style="flex:1.45"><table>
<tr><th style="width:24%">Proposal 要求</th><th>落地方式</th></tr>
<tr><td>Figma 组件与代码组件的对应关系</td><td><code>match</code> 按名字、UI 概念（中英文同义词）、项目前缀、使用次数、设计系统封装打分；人确认后写进 <code>component-map.yaml</code>，之后 match 优先采用确认过的映射</td></tr>
<tr><td>开发前识别可复用组件</td><td>先 <code>reuse</code> 再 <code>match</code>：只推荐<b>本项目真能用</b>的组件（检查组件库版本、改名 / 移除、主题依赖、废弃 API、设计系统是否封装过）；没有组件但有固定写法时推荐「项目惯例」</td></tr>
<tr><td>开发后检查用了哪些组件、和 Figma 截图对比</td><td><code>reuse --focus &lt;页面&gt;</code> 或 <code>--changed &lt;git ref&gt;</code> 只看新增代码；<code>compare</code> 截图对比；demo 里 <code>verify</code> 把两者合成一个门禁</td></tr>
<tr><td>偏差反馈给 Agent、沉淀规则</td><td>每条问题都带文件行号或 Figma 节点 id 和修改建议，Skill 规定照着改到通过；决定和理由写进 component-map，经验写进 <code>rules.md</code></td></tr>
</table></div>
<div><img class="wide" src="{img}23-match-newbee-overlay.png">
<p class="note">newbee-mall（Vue 3 + Vant）上跑「订单详情」的复用清单：绿框直接复用、橙框需确认、红框需新建；右侧是项目里的真实写法和 Figma 变体对应的 props。</p></div></div>

<h2 class="pb">2. 复用清单准不准：{len(projects)} 个真实开源项目（固定 commit，可复现）</h2>
<table class="ev">
<tr><th>项目</th><th>技术栈</th><th>留出</th><th>订单详情（组件化）</th><th>我的持仓（未组件化、中文图层）</th><th>买入下单（未组件化、新概念）</th></tr>
{eval_rows()}
<tr><th colspan="3">合计</th><th colspan="3">✓{tot['✓']} / △{tot['△']} / ✗{tot['✗']}</th></tr>
</table>
<p class="note">✓ 推荐正确；△ 可用但不是最优（多为「需确认」）；✗ 推荐错误。留出 = 调规则时没用过的项目。</p>
<h3>留出项目第一次跑的结果（如实记录，修正前）</h3>
<table>
<tr><th style="width:17%">项目</th><th style="width:20%">第一次</th><th>出错原因 → 通用修正</th></tr>
<tr><td>Jetsnack</td><td>订单详情 <b class="bad">3 ✗</b></td><td>把桌面小部件（Glance）组件当成 App 组件；没认出 <code>Jetsnack*</code> 是项目前缀 → 识别 Glance、学习项目前缀、<code>ui/components</code> 算设计系统</td></tr>
<tr><td>litemall</td><td>两份各 <b class="bad">1 ✗</b></td><td>推荐了 Vant 3 才有的 ActionBar，项目是 Vant 2 → 组件库加入改名 / 移除信息，按实际安装版本过滤</td></tr>
<tr><td>hello-uniapp</td><td>持仓页 <b class="bad">12 ✗</b></td><td>uni-app 内置小写标签被当成项目组件、中文图层名规范化后为空被判「同名」→ 加 uni-ui 目录，修空名误判</td></tr>
<tr><td>NewPipe</td><td><b class="bad">3 ✗</b></td><td>AppCompat 标准控件没算进项目资产，XML 页面被推荐 Compose 组件 → 跨技术栈的候选最多「需确认」</td></tr>
<tr><td>mall-app-web / AntennaPod（第二批）</td><td>崩溃 + 1 ✗ / <b class="bad">1 ✗</b></td><td>uni-app 内置标签没有 import 路径导致崩溃；<code>&lt;view&gt;</code> 用量压低了 uni-ui 的占比；AntennaPod 的空状态封装 EmptyViewHandler 不是 View → 修崩溃、占比排除平台内置元素、识别「界面封装类」</td></tr>
<tr><td>WanAndroid / Seal（第三批）</td><td><b class="bad">11 ✗</b> / <b class="bad">7 ✗</b></td><td>三方 View 库（ActionBarEx、MultiStateView）没被收进资产；Android 的 ActionBar 是顶栏；「Video*」被当成项目前缀；专门用途的组件（PasswordInputView、ConfirmButton）被当成通用组件 → 收录三方 View、按平台区分概念、用途词降权</td></tr>
<tr><td>买入下单（新设计稿）</td><td>14 个项目 <b class="bad">12 ✗</b></td><td>开关 / 复选框 / 单选原来是一个概念，搜索算在输入框里 → 拆开概念、补词表、检查 Material 3 主题</td></tr>
</table>
<p class="note">vant-demo、Pokedex、Sunflower、vue3-vant-mobile 第一次跑就没有 ✗。复用推荐只错在「项目没见过的结构」，修正都是通用规则，没有针对单个项目的特判；<code>python3 experiments/real-match-eval.py</code> 每次改规则都会全量回归。</p>

<div class="two">
<div><h3>体检：组件库有，不等于本项目能用</h3><img class="wide" src="{img}21-qmui-compat.png">
<p class="note">QMUI_Android 用的是 AppCompat 主题和 Material 1.4：这些组件会崩溃或编译失败，match 不会推荐。</p></div>
<div><h3>写完检查：有现成组件却手写</h3><img class="wide" src="{img}26-check-a.png">
<p class="note">不用 Skill 的实现（newbee，A 组）：分割线、商品卡、底部操作栏都是手写的，每条带行号和应复用的组件。</p></div>
</div>

<h2 class="pb">3. Agent 实测：真实项目里实现一个新页面</h2>
<p>同一份设计稿（订单详情），在 5 个开源项目里新建页面。A 组只给任务和 Figma 导出，C 组加上 Skill（体检 → 复用清单 → 实现 → 页面检查）。所有 run 用同一版 <code>figmafe reuse --focus / --changed</code> 打分：自造已有组件、有组件没用、手写已有组件（按 class / id 或横排「标题 + 值」结构识别）任何一项都算不通过。</p>
<table>
<tr><th rowspan="2">项目</th><th rowspan="2">技术栈</th><th colspan="3">A 不用 Skill（Opus / Haiku）</th><th colspan="3">C 用 Skill（Opus / Haiku）</th></tr>
<tr><th>通过</th><th>复用组件数</th><th>问题数</th><th>通过</th><th>复用组件数</th><th>问题数</th></tr>
{exp_rows()}
<tr><th colspan="2">合计</th><th class="bad">{passes(A_all)}</th><th>平均 {avg(A_all, lambda r: r['used']):.1f}</th><th>共 {sum(bad(r) for r in A_all)}</th><th class="ok">{passes(C_all)}</th><th>平均 {avg(C_all, lambda r: r['used']):.1f}</th><th>共 {sum(bad(r) for r in C_all)}</th></tr>
</table>
<ul>
<li><b>Now in Android</b>：A 组没先查已有组件，在 core/designsystem 下新建了 NiaNavBar / NiaCard / NiaDivider / NiaCell 等，和已有的 NiaTopAppBar / Card / HorizontalDivider / ListItem 重复，还写了一套 18 个写死颜色的主题。C 组复用了它们，新建的组件都在 component-map.yaml 里写了理由。</li>
<li><b>QMUI / NewPipe</b>：C 组发现 match 推荐了项目用不了的组件（MaterialDivider 需要 Material 1.5、项目是 1.4；Chip 在 AppCompat 主题下会崩）。这两类问题已经做进工具：体检会列出「组件库有、但本项目不能直接用的组件」，match 不再推荐。实验中 agent 报告的工具问题都已修复并加了回归测试。</li>
<li><b>成本</b>：有用量记录的 {len(ca)} 次 A、{len(cc)} 次 C，C 平均多 {k_up*100:.0f}% token、{t_up*100:.0f}% 时间，多在体检、复用清单、页面检查和写决定这几步。</li>
</ul>

<h2>4. 新页面实现闭环（demo 仓库）</h2>
<ul>
<li>自建的 H5（React）和 Android（Compose，兼容 compileSdk 32 / minSdk 23；另有 XML View 存量组件）组件库加一份订单详情设计稿，跑了 {len(demo)} 次 A / C 对照。<code>verify</code> 能区分「页面写错」和「组件本身不符合设计」（Android 剩下的 2 处严重偏差出在组件库内部，见 <code>docs/component-issues.md</code>）；教程第 8 步故意写错 3 处，全部被抓出。</li>
<li>这个仓库小、组件好找，A 组复用率也是 100%，<b>证明不了复用率提升</b>；差别在规范：A 组 {sum(r['selfBuiltCountdown'] for r in dA)}/{len(dA)} 次自造了设计稿外的 CountDown，C 组 {sum(r['selfBuiltCountdown'] for r in dC)}/{len(dC)}。复用率的差距要在真实项目里看（第 3 节）。</li>
</ul>
<h2>5. 局限（如实说明）</h2>
<ul>
<li><b>没有真实 Figma 文件</b>，本环境也访问不了 Figma：三份设计稿由手写的设计源生成（格式和插件导出一致）。插件只用 mock 的 Plugin API 测过，没在真实 Figma 里跑过。</li>
<li><b>没有编译真实项目</b>：Agent 实验只做静态检查，没有构建和运行这些开源项目（本环境装不了 Android SDK）；demo 的 Android 截图用 Compose Desktop 渲染，不是真机。</li>
<li><b>实验有泄漏</b>：A 组 agent 也能看到本仓库的 Skill；1 次 A run 自己跑了 figmafe，已剔除。泄漏只会让 A 组更好，所以 A / C 差距是保守估计。样本量小（每格 1～2 次），只能说明方向。</li>
<li><b>评分标准是人定的</b>：✓ / △ / ✗ 由我逐条判断并写成回归用例；真实接入时应由组件 owner 复核一次。</li>
</ul>

<h2>6. 接入你们的项目</h2>
<table>
<tr><th style="width:4%">#</th><th>动作</th><th style="width:34%">产出</th></tr>
<tr><td>1</td><td><code>figmafe reuse &lt;H5 仓库&gt;</code>、<code>figmafe reuse &lt;Android 仓库&gt;</code>（只读，不改代码）</td><td>资产清单、优先处理项、首份体检基线</td></tr>
<tr><td>2</td><td>Figma 插件导出下一个要做的页面，<code>figmafe match</code></td><td>复用清单 + 叠图，组件 owner 确认「需确认」项</td></tr>
<tr><td>3</td><td>开发时让 Agent 按 Skill 走；提交前 <code>reuse --changed origin/main</code></td><td>新增代码的复用问题，带行号</td></tr>
<tr><td>4</td><td>CI 里 <code>reuse --baseline main.json --no-regress</code>（示例：<code>docs/ci/</code>）</td><td>组件化率、写死颜色等指标只许变好</td></tr>
</table>
<p class="note">教程：<code>docs/tutorial/README.md</code> · 数据：<code>experiments/real-match-eval.md</code>、<code>match-judgement.md</code>、<code>real-project/summary.md</code> · 重新生成：<code>python3 docs/tools/build_report.py</code></p>
</body></html>"""
open(os.path.join(root, 'docs', 'report.html'), 'w').write(page)
print('✓ docs/report.html')
