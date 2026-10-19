#!/usr/bin/env python3
"""Summarise the real brownfield runs (experiments/real-project/<run>/<platform>.page.json) into summary.md / summary.json."""
import json, os, glob

HERE = os.path.dirname(os.path.abspath(__file__))
RP = os.path.join(HERE, "real-project")

PROJECTS = {
    "Rnb": ("newbee-mall-vue3-app（Vue 3 + Vant）", "h5", "新建 `src/views/OrderDetailV2.vue`；打分：`figmafe reuse <项目> --focus src/views/OrderDetailV2.vue`"),
    "Rqm": ("QMUI_Android（XML View + QMUI）", "android", "新建 Fragment + layout；打分：`figmafe reuse <项目> --changed HEAD`（新文件整份判，已有文件只看新增行）"),
    "Rnia": ("Now in Android（Jetpack Compose + Material3 + 自有设计系统 Nia*）", "android", "新建 `OrderDetailScreen.kt`，通用组件放哪由 agent 决定；打分：`figmafe reuse <项目> --changed HEAD`"),
}
GROUP = {"A": "A 只给任务（Opus）", "Ahk": "A 只给任务（Haiku）", "C": "C 用 Skill（Opus）", "Chk": "C 用 Skill（Haiku）"}
# what each agent created as new components (from their final reports)
NEW = {
    "Rnb-C-1": "Price.vue", "Rnb-C-2": "Price.vue", "Rnb-Chk-1": "Price.vue",
    "Rqm-C-1": "QDNoticeBarView、QDCountDownView、QDPriceView、QDActionBar",
    "Rqm-C-2": "QDNoticeBarView、QDCountDownView、QDPriceView、QDActionBarLayout、QDTagView、QDProductCardView",
    "Rnia-A-1": "12 个：NiaNavBar、NiaCard、NiaDivider、NiaCell、NiaPillButton、NiaStatusTag、NiaPrice、NiaCountDown、NiaNoticeBar、NiaActionBar、NiaProductCard + CommerceTheme",
    "Rnia-C-1": "NiaNoticeBar、NiaPriceText、NiaCountdownText、NiaActionBar、NiaListItem、NiaStatusTag、ProductCard",
}

rows, out = [], {}
for pj in sorted(glob.glob(os.path.join(RP, "R*", "*.page.json"))):
    run = os.path.basename(os.path.dirname(pj))
    proj, grp = run.split("-")[0], run.split("-")[1]
    d = json.load(open(pj))
    err = lambda fs: sum(1 for f in fs or [] if f.get("severity") in ("error", "warn") and f.get("severity") != "info")
    hb = d.get("handBuilt") or []
    rec = {
        "run": run, "project": proj, "group": grp, "pass": d["pass"], "used": len(d.get("used") or []),
        "reinvented": err(d.get("reinvented")), "missed": err(d.get("missedReuse")),
        "handBuiltError": sum(1 for f in hb if f["severity"] == "error"), "handBuiltWarn": sum(1 for f in hb if f["severity"] != "error"),
        "localTokens": d.get("localTokens", 0), "colors": len(d.get("colors") or []), "new": NEW.get(run, "—"),
    }
    out[run] = rec
    rows.append(rec)

md = ["# 真实存量项目实验：拿到新设计稿「订单详情」，在开源项目里实现", "",
      "A = 只给任务和 Figma 导出；C = 加上 Skill（`reuse` 体检 → `match` 复用清单 → 实现 → 页面复用检查）。所有 run 用同一版 `figmafe` 打分，由 `experiments/real-summary.py` 生成。", ""]
for p, (name, plat, how) in PROJECTS.items():
    rs = [r for r in rows if r["project"] == p]
    if not rs:
        continue
    md += [f"## {name}", "", how, "",
           "| run | 组 | 结果 | 复用现成组件 | 自造轮子 | 有组件没用 | 手写已有组件（确定/疑似） | 新增颜色变量 | 写死颜色 | 新建组件 |",
           "|---|---|---|---|---|---|---|---|---|---|"]
    for r in rs:
        md.append(f"| {r['run']} | {GROUP.get(r['group'], r['group'])} | {'✅ PASS' if r['pass'] else '❌ FAIL'} | {r['used']} | {r['reinvented']} | {r['missed']} | {r['handBuiltError']} / {r['handBuiltWarn']} | {r['localTokens']} | {r['colors']} | {r['new']} |")
    md.append("")
md += ["## 说明", "",
       "- newbee：A 组页面里手写了商品卡、分割线、底部操作栏，还自建了 10 个颜色变量；C 组直接用 Vant 的 Card / Divider / ActionBar 和项目已有的 SimpleHeader，缺的 Price 放进组件目录。",
       "- Haiku 不带 Skill 时一个现成组件都没用（0 个），带 Skill 后复用 10 个、检查通过：Skill 把「先查再写」变成了固定流程，模型能力弱时收益更大。",
       "- QMUI：A 组地址行用 LinearLayout 手写（项目里有 QMUICommonListItemView，已用 48 次），C 组用 QMUIGroupListView + QMUICommonListItemView。",
       "- QMUI 的两个 C run 都发现 match 推荐了项目用不了的组件（MaterialDivider 需要 Material 1.5，项目是 1.4；Chip / MaterialCardView 在 AppCompat 主题下会崩溃），自行绕开并写进了报告。这两类问题已经做进工具：现在体检会列出「组件库里有、但本项目不能直接用的组件」，match 不再推荐它们。",
       "- Now in Android：A 组没有先查项目已有组件，在 core/designsystem 下新建了 12 个组件，其中 NiaNavBar / NiaCard / NiaDivider / NiaPillButton / NiaCell 和已有的 NiaTopAppBar / Card / HorizontalDivider / NiaButton / ListItem 重复，另外写了一套 18 个写死颜色的主题。C 组复用了 NiaTopAppBar、Card、HorizontalDivider、NiaButton、NiaOutlinedButton，颜色全部映射到 MaterialTheme；它按语义拒绝了 match 的两个 △ 推荐（InterestsItem、NiaFilterChip），新建的 NiaListItem 被 PR 闸门判为和 Material3 ListItem 重复，属于需要人定的项，agent 在交付里如实说明了理由。",
       "- C 组 agent 会把这次学到的经验追加到 rules.md（例如「直接复用也要核对组件语义」），这正是方案里「沉淀规则」那一步；其中有价值的已合并进 Skill。",
       "- 新增颜色变量：设计稿里的颜色项目里大多没有同值 token，C 组把它们集中加到 colors.xml；检查只把「与已有 token 同值却另起名字」的算问题。",
       ""]
open(os.path.join(RP, "summary.md"), "w").write("\n".join(md))
json.dump(out, open(os.path.join(RP, "summary.json"), "w"), ensure_ascii=False, indent=1)
print("\n".join(md))
