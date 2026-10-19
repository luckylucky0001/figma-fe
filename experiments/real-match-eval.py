#!/usr/bin/env python3
"""Regression check for `figmafe match` on real open-source projects.

Clones each project at a pinned commit into experiments/.cache/, runs match for both Figma exports
and grades the top recommendation against hand-judged answers (experiments/match-judgement.md):
✓ = one of the right answers, △ = usable but not the best, ✗ = wrong. "—" means "needs a new component".

    python3 experiments/real-match-eval.py            # all projects
    python3 experiments/real-match-eval.py nia        # one project
"""
import json, os, subprocess, sys, tempfile

HERE = os.path.dirname(os.path.abspath(__file__))
ROOT = os.path.dirname(HERE)
CACHE = os.path.join(HERE, ".cache")
BIN = os.path.join(ROOT, "bin", "figmafe")

PROJECTS = {
    "newbee": ("https://github.com/newbee-ltd/newbee-mall-vue3-app.git", "07740cba11c219d971879282a9faa35ca9e4cd64", "h5"),
    "qmui": ("https://github.com/Tencent/QMUI_Android.git", "026e7d486677d6593a96689cd590ec3561176717", "android"),
    "nia": ("https://github.com/android/nowinandroid.git", "a49ed253d75e61a2b6ab80a8da677b57437b08eb", "android"),
    "antdm-template": ("https://github.com/lhw-dev/react-h5-template.git", "d0fc97327b35a2675f55a2fe88c5b39d20aec8a1", "h5"),
    "antdm-mobile": ("https://github.com/xmllein/react-mobile.git", "12e4276ce83962d92a86b11fa2e89c1ec94f927f", "h5"),
    # held out: first run after the rules above were tuned (see match-judgement.md for the first-pass result)
    "jetsnack": ("https://github.com/android/compose-samples.git#Jetsnack", "0bc770fd28510b5a16395349e7281f11386ac82f", "android"),
    "newpipe": ("https://github.com/TeamNewPipe/NewPipe.git", "32eb85d0fb09ef9022421751fd96bee564b9223b", "android"),
    "vant-demo-vue3": ("https://github.com/youzan/vant-demo.git#vant/vue3-ts", "c382aeba3cfea6ea31f6fee98a99b1c6604accc4", "h5"),
    "vant-demo-vue2": ("https://github.com/youzan/vant-demo.git#vant/vue2", "c382aeba3cfea6ea31f6fee98a99b1c6604accc4", "h5"),
}

# figma element → (right answers, acceptable answers)
ANTD = {
    "NavBar": (["NavBar"], []), "NoticeBar": (["NoticeBar"], []), "Card": (["Card"], []), "Tag": (["Tag"], []),
    "CountDown": (["—"], []), "Divider": (["Divider"], []), "Cell": (["ListItem"], []), "ProductCard": (["Card", "—"], []),
    "Price": (["—"], []), "ActionBar": (["—"], []), "Button": (["Button"], []),
}
ANTD_H = {
    "导航栏": (["NavBar"], []), "公告通知栏": (["NoticeBar"], []), "资产卡片": (["Card"], []), "标签页": (["Tabs"], []),
    "持仓列表": (["List"], []), "持仓列表项": (["ListItem"], []), "涨跌标签": (["Tag"], []), "收盘倒计时": (["—"], []),
    "底部操作栏": (["—"], []), "买入按钮": (["Button"], []), "卖出按钮": (["Button"], []), "更多按钮": (["Button"], []),
}
VANT = {
    "NavBar": (["NavBar"], []), "NoticeBar": (["NoticeBar"], []), "Card": (["CellGroup"], ["—"]), "Tag": (["Tag"], []),
    "CountDown": (["CountDown"], []), "Divider": (["Divider"], []), "Cell": (["Cell"], []), "ProductCard": (["Card"], []),
    "Price": (["—"], []), "ActionBar": (["ActionBar"], []), "Button": (["Button"], []),
}
VANT_H = {
    "导航栏": (["NavBar"], []), "公告通知栏": (["NoticeBar"], []), "资产卡片": (["CellGroup", "—"], []), "标签页": (["Tabs"], []),
    "持仓列表": (["List"], []), "持仓列表项": (["Cell"], []), "涨跌标签": (["Tag"], []), "收盘倒计时": (["CountDown"], []),
    "底部操作栏": (["ActionBar"], []), "买入按钮": (["Button"], []), "卖出按钮": (["Button"], []), "更多按钮": (["Button"], []),
}
EXPECTED = {
    ("newbee", "order-detail"): {
        "NavBar": (["SimpleHeader"], ["NavBar"]), "NoticeBar": (["NoticeBar"], []), "Card": (["CellGroup"], []), "Tag": (["Tag"], []),
        "CountDown": (["CountDown"], []), "Divider": (["Divider", "Vant 细线类 van-hairline--*"], []), "Cell": (["Cell"], []),
        "ProductCard": (["Card"], []), "Price": (["—"], []), "ActionBar": (["ActionBar"], []), "Button": (["Button"], []),
    },
    ("newbee", "holdings"): {
        "导航栏": (["SimpleHeader"], ["NavBar"]), "公告通知栏": (["NoticeBar"], []), "资产卡片": (["CellGroup", "—"], []), "标签页": (["Tabs"], []),
        "持仓列表": (["List"], []), "持仓列表项": (["Cell"], []), "涨跌标签": (["Tag"], []), "收盘倒计时": (["CountDown"], []),
        "底部操作栏": (["ActionBar"], []), "买入按钮": (["Button"], []), "卖出按钮": (["Button"], []), "更多按钮": (["Button"], []),
    },
    ("qmui", "order-detail"): {
        "NavBar": (["QMUITopBarLayout", "QMUITopBar"], []), "NoticeBar": (["—"], []), "Card": (["QMUIGroupListView"], []), "Tag": (["—"], ["QMUIRoundButton"]),
        "CountDown": (["—"], []), "Divider": (["QMUILinearLayout", "QMUIFrameLayout"], ["QMUIPriorityLinearLayout", "QMUIConstraintLayout", "—"]),
        "Cell": (["QMUICommonListItemView"], []), "ProductCard": (["—"], []), "Price": (["—"], []), "ActionBar": (["—"], []), "Button": (["QMUIRoundButton"], ["QMUIButton"]),
    },
    ("qmui", "holdings"): {
        "导航栏": (["QMUITopBarLayout", "QMUITopBar"], []), "公告通知栏": (["—"], []), "资产卡片": (["QMUIGroupListView", "—"], []), "标签页": (["QMUITabSegment"], ["QMUIBasicTabSegment"]),
        "持仓列表": (["RecyclerView"], []), "持仓列表项": (["QMUICommonListItemView"], []), "涨跌标签": (["—"], ["QMUIRoundButton"]), "收盘倒计时": (["—"], []),
        "底部操作栏": (["—"], []), "买入按钮": (["QMUIRoundButton"], ["QMUIButton"]), "卖出按钮": (["QMUIRoundButton"], ["QMUIButton"]), "更多按钮": (["QMUIRoundButton"], ["QMUIButton"]),
    },
    ("nia", "order-detail"): {
        "NavBar": (["NiaTopAppBar"], []), "NoticeBar": (["—"], []), "Card": (["Card"], []), "Tag": (["NiaTopicTag"], ["NiaFilterChip"]),
        "CountDown": (["—"], []), "Divider": (["HorizontalDivider"], []), "Cell": (["ListItem"], ["InterestsItem"]), "ProductCard": (["Card"], ["—"]),
        "Price": (["—"], []), "ActionBar": (["—"], []), "Button": (["NiaButton"], ["NiaTextButton", "NiaOutlinedButton"]),
    },
    ("nia", "holdings"): {
        "导航栏": (["NiaTopAppBar"], []), "公告通知栏": (["—"], []), "资产卡片": (["Card"], []), "标签页": (["NiaTabRow"], ["NiaTab"]),
        "持仓列表": (["LazyColumn"], ["—"]), "持仓列表项": (["ListItem"], ["InterestsItem"]), "涨跌标签": (["NiaTopicTag"], ["NiaFilterChip"]), "收盘倒计时": (["—"], []),
        "底部操作栏": (["—"], []), "买入按钮": (["NiaButton"], []), "卖出按钮": (["NiaButton"], []), "更多按钮": (["NiaButton"], ["NiaOutlinedButton"]),
    },
    ("jetsnack", "order-detail"): {
        "NavBar": (["TopAppBar", "CenterAlignedTopAppBar"], ["—"]), "NoticeBar": (["—"], []), "Card": (["JetsnackCard"], ["Card"]), "Tag": (["—", "Badge"], ["FilterChip"]),
        "CountDown": (["—"], []), "Divider": (["JetsnackDivider"], ["HorizontalDivider"]), "Cell": (["ListItem"], ["—"]), "ProductCard": (["JetsnackCard", "—"], ["Card"]),
        "Price": (["—"], []), "ActionBar": (["—"], []), "Button": (["JetsnackButton"], []),
    },
    ("jetsnack", "holdings"): {
        "导航栏": (["TopAppBar", "CenterAlignedTopAppBar"], ["—"]), "公告通知栏": (["—"], []), "资产卡片": (["JetsnackCard"], ["Card"]), "标签页": (["TabRow"], ["—"]),
        "持仓列表": (["LazyColumn"], []), "持仓列表项": (["ListItem"], ["—"]), "涨跌标签": (["—", "Badge"], ["FilterChip"]), "收盘倒计时": (["—"], []),
        "底部操作栏": (["—"], []), "买入按钮": (["JetsnackButton"], []), "卖出按钮": (["JetsnackButton"], []), "更多按钮": (["JetsnackButton"], []),
    },
    ("newpipe", "order-detail"): {
        "NavBar": (["Toolbar", "MaterialToolbar"], []), "NoticeBar": (["—"], []), "Card": (["CardView", "MaterialCardView"], []), "Tag": (["Chip"], ["—"]),
        "CountDown": (["—"], []), "Divider": (["MaterialDivider", "—"], []), "Cell": (["—"], ["PreferenceRow", "ListItem"]), "ProductCard": (["CardView", "—"], []),
        "Price": (["—"], []), "ActionBar": (["—"], []), "Button": (["MaterialButton"], ["Button"]),
    },
    ("newpipe", "holdings"): {
        "导航栏": (["Toolbar", "MaterialToolbar"], []), "公告通知栏": (["—"], []), "资产卡片": (["CardView", "MaterialCardView"], []), "标签页": (["TabLayout"], []),
        "持仓列表": (["RecyclerView"], []), "持仓列表项": (["—"], ["PreferenceRow", "ListItem"]), "涨跌标签": (["Chip"], ["—"]), "收盘倒计时": (["—"], []),
        "底部操作栏": (["—"], []), "买入按钮": (["MaterialButton"], ["Button"]), "卖出按钮": (["MaterialButton"], ["Button"]), "更多按钮": (["MaterialButton"], ["Button"]),
    },
    ("vant-demo-vue3", "order-detail"): "VANT", ("vant-demo-vue3", "holdings"): "VANT_H",
    ("vant-demo-vue2", "order-detail"): "VANT", ("vant-demo-vue2", "holdings"): "VANT_H",
    ("antdm-template", "order-detail"): ANTD, ("antdm-template", "holdings"): ANTD_H,
    ("antdm-mobile", "order-detail"): ANTD, ("antdm-mobile", "holdings"): ANTD_H,
}


def checkout(name, url, commit):
    """url may end with #subdir: only that directory is checked out (sparse) and returned."""
    url, _, sub = url.partition("#")
    d = os.path.join(CACHE, name)
    if not os.path.isdir(os.path.join(d, ".git")):
        os.makedirs(d, exist_ok=True)
        subprocess.run(["git", "init", "-q", d], check=True)
        subprocess.run(["git", "-C", d, "remote", "add", "origin", url], check=True)
        if sub:
            subprocess.run(["git", "-C", d, "sparse-checkout", "set", sub], check=True)
    have = subprocess.run(["git", "-C", d, "rev-parse", "HEAD"], capture_output=True, text=True).stdout.strip()
    if have != commit:
        subprocess.run(["git", "-C", d, "fetch", "-q", "--depth", "1", "--filter=blob:none", "origin", commit], check=True)
        subprocess.run(["git", "-C", d, "checkout", "-q", commit], check=True)
    return os.path.join(d, sub) if sub else d


def main():
    only = sys.argv[1:]
    total = {"✓": 0, "△": 0, "✗": 0}
    report = {}
    for name, (url, commit, pf) in PROJECTS.items():
        if only and name not in only:
            continue
        d = checkout(name, url, commit)
        for page in ("order-detail", "holdings"):
            exp = EXPECTED.get((name, page))
            if isinstance(exp, str):
                exp = globals()[exp]
            if not exp:
                continue
            with tempfile.TemporaryDirectory() as out:
                subprocess.run([BIN, "match", "--figma", os.path.join(ROOT, "design/figma-export", page), d, "--platform", pf, "-o", out, "-q"], check=True, capture_output=True)
                m = json.load(open(os.path.join(out, pf + ".json")))
            counts = {"✓": 0, "△": 0, "✗": 0}
            lines = []
            for it in m["items"]:
                key = it["figma"].replace("（未组件化）", "")
                top = it["candidates"][0]["name"] if it["candidates"] and it["verdict"] != "new" else "—"
                right, ok = exp.get(key, ([], []))
                g = "✓" if top in right else "△" if top in ok else "✗"
                counts[g] += 1
                total[g] += 1
                if g != "✓":
                    lines.append(f"    {g} {key}: {top}（期望 {' / '.join(right)}）")
            report[f"{name}/{page}"] = counts
            print(f"{name:15s} {page:13s} ✓{counts['✓']:2d} △{counts['△']} ✗{counts['✗']}")
            print("\n".join(lines)) if lines else None
    print(f"合计 ✓{total['✓']} △{total['△']} ✗{total['✗']}")
    json.dump(report, open(os.path.join(HERE, "real-match-eval.json"), "w"), ensure_ascii=False, indent=1)
    sys.exit(1 if total["✗"] else 0)


if __name__ == "__main__":
    main()
