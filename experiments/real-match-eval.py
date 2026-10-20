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
    "litemall": ("https://github.com/linlinjava/litemall.git#litemall-vue", "a1ef964a718b7277925b19ea26afe78ea3a1d325", "h5"),
    "hello-uniapp": ("https://github.com/dcloudio/hello-uniapp.git", "ab80c34681013d844be2a9feb362bf96ffeca3ee", "h5"),
    "taro-ui-demo": ("https://github.com/NervJS/taro-ui-demo.git", "8d87bb2e3da32fe5577b1f7daebaa6261b5b3f72", "h5"),
    "nutui-vite": ("https://github.com/jdf2e/nutui-demo.git#vite", "e84f77511f2f5b81eba56763829990dfd1a579b8", "h5"),
    "sunflower": ("https://github.com/android/sunflower.git", "2a357a31551bb53f3fe80382a9ce6d30bcc8b960", "android"),
    "pokedex": ("https://github.com/skydoves/Pokedex.git", "769eb47af4a4a98222253a67026ba42b99e2fc32", "android"),
    "newpipe": ("https://github.com/TeamNewPipe/NewPipe.git", "32eb85d0fb09ef9022421751fd96bee564b9223b", "android"),
    "vant-demo-vue3": ("https://github.com/youzan/vant-demo.git#vant/vue3-ts", "c382aeba3cfea6ea31f6fee98a99b1c6604accc4", "h5"),
    "vant-demo-vue2": ("https://github.com/youzan/vant-demo.git#vant/vue2", "c382aeba3cfea6ea31f6fee98a99b1c6604accc4", "h5"),
    # held out, second batch (added after the report's first version)
    "vue3-vant-mobile": ("https://github.com/easy-temps/vue3-vant-mobile.git", "8a46a5a41473ba3a638e2249b4c4e48b8ef7ddea", "h5"),
    "mall-app-web": ("https://github.com/macrozheng/mall-app-web.git", "b1eb3cde1d74fbf62e72195bafd9322c8149f4c8", "h5"),
    "antennapod": ("https://github.com/AntennaPod/AntennaPod.git", "3d105d598dfa582c6cb22ba40b56d11780366890", "android"),
    "wanandroid": ("https://github.com/goweii/WanAndroid.git", "baba04a90460be2233ddd715fa41f1faa9be741d", "android"),
    "seal": ("https://github.com/JunkFood02/Seal.git", "7677f61a20fda4210e84215fef9e9b51d25daa47", "android"),
    "kunminx": ("https://github.com/KunMinX/Jetpack-MVVM-Best-Practice.git", "543eb8659089d74ccad403763cb16596febc89b7", "android"),
    "jeecg-uniapp": ("https://github.com/jeecgboot/jeecg-uniapp.git", "1f33e0225f6557994805db6706034cc0db0232c7", "h5"),
    "vue-h5-template": ("https://github.com/sunniejs/vue-h5-template.git", "9ac24512ea0eb36cf35f54803e9bf267af18335d", "h5"),
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
# Vant 2 has no ActionBar: its action bar is GoodsAction (renamed in Vant 3)
VANT2 = dict(VANT, ActionBar=(["GoodsAction"], []))
VANT2_H = dict(VANT_H, 底部操作栏=(["GoodsAction"], []))
# litemall registers @chenfengyuan/vue-countdown globally and uses it on 5 pages: that is its countdown
LITEMALL = dict(VANT2, CountDown=(["VueCountdown"], ["CountDown"]))
LITEMALL_H = dict(VANT2_H, 收盘倒计时=(["VueCountdown"], ["CountDown"]))
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
    ("hello-uniapp", "order-detail"): {
        "NavBar": (["NavBar"], []), "NoticeBar": (["NoticeBar"], []), "Card": (["Card"], []), "Tag": (["Tag"], []),
        # .uni-divider is the divider style in the project's own common/uni.css (accepted after the run that found it)
        "CountDown": (["Countdown"], []), "Divider": (["—", "CSS 类 .uni-divider"], ["Divider"]), "Cell": (["ListItem"], []), "ProductCard": (["Card", "—"], []),
        "Price": (["—"], []), "ActionBar": (["GoodsNav"], []), "Button": (["Button[uni-app 内置组件]"], []),
    },
    ("hello-uniapp", "holdings"): {
        "导航栏": (["NavBar"], []), "公告通知栏": (["NoticeBar"], []), "资产卡片": (["Card"], []), "标签页": (["SegmentedControl"], []),
        "持仓列表": (["List"], []), "持仓列表项": (["ListItem"], []), "涨跌标签": (["Tag"], []), "收盘倒计时": (["Countdown"], []),
        "底部操作栏": (["GoodsNav"], []), "买入按钮": (["Button[uni-app 内置组件]"], []), "卖出按钮": (["Button[uni-app 内置组件]"], []), "更多按钮": (["Button[uni-app 内置组件]"], []),
    },
    ("taro-ui-demo", "order-detail"): {
        "NavBar": (["AtNavBar"], []), "NoticeBar": (["AtNoticebar"], []), "Card": (["AtCard"], []), "Tag": (["AtTag"], []),
        "CountDown": (["AtCountdown"], []), "Divider": (["AtDivider"], []), "Cell": (["AtListItem"], []), "ProductCard": (["—", "AtCard"], []),
        "Price": (["—"], []), "ActionBar": (["—"], []), "Button": (["AtButton"], []),
    },
    ("taro-ui-demo", "holdings"): {
        "导航栏": (["AtNavBar"], []), "公告通知栏": (["AtNoticebar"], []), "资产卡片": (["AtCard"], []), "标签页": (["AtTabs", "AtSegmentedControl"], []),
        "持仓列表": (["AtList", "—"], []), "持仓列表项": (["AtListItem"], []), "涨跌标签": (["AtTag"], []), "收盘倒计时": (["AtCountdown"], []),
        "底部操作栏": (["—"], []), "买入按钮": (["AtButton"], []), "卖出按钮": (["AtButton"], []), "更多按钮": (["AtButton"], []),
    },
    ("nutui-vite", "order-detail"): {
        "NavBar": (["Navbar"], []), "NoticeBar": (["NoticeBar"], []), "Card": (["CellGroup"], ["—"]), "Tag": (["Tag"], []),
        "CountDown": (["Countdown"], []), "Divider": (["Divider"], []), "Cell": (["Cell"], []), "ProductCard": (["Card"], ["—"]),
        "Price": (["Price"], []), "ActionBar": (["—"], []), "Button": (["Button"], []),
    },
    ("nutui-vite", "holdings"): {
        "导航栏": (["Navbar"], []), "公告通知栏": (["NoticeBar"], []), "资产卡片": (["CellGroup"], ["—"]), "标签页": (["Tabs"], []),
        "持仓列表": (["List", "—"], []), "持仓列表项": (["Cell"], []), "涨跌标签": (["Tag"], []), "收盘倒计时": (["Countdown"], []),
        "底部操作栏": (["—"], []), "买入按钮": (["Button"], []), "卖出按钮": (["Button"], []), "更多按钮": (["Button"], []),
    },
    ("pokedex", "order-detail"): {
        "NavBar": (["MaterialToolbar", "Toolbar"], []), "NoticeBar": (["—"], []), "Card": (["MaterialCardView", "CardView"], []), "Tag": (["Chip"], []),
        "CountDown": (["—"], []), "Divider": (["MaterialDivider"], []), "Cell": (["—"], []), "ProductCard": (["MaterialCardView", "—"], []),
        "Price": (["—"], []), "ActionBar": (["—"], []), "Button": (["MaterialButton"], []),
    },
    ("pokedex", "holdings"): {
        "导航栏": (["MaterialToolbar", "Toolbar"], []), "公告通知栏": (["—"], []), "资产卡片": (["MaterialCardView", "CardView"], []), "标签页": (["TabLayout"], []),
        "持仓列表": (["RecyclerView"], []), "持仓列表项": (["—"], []), "涨跌标签": (["Chip"], []), "收盘倒计时": (["—"], []),
        "底部操作栏": (["—"], []), "买入按钮": (["MaterialButton"], []), "卖出按钮": (["MaterialButton"], []), "更多按钮": (["MaterialButton"], []),
    },
    ("newpipe", "order-detail"): {
        "NavBar": (["Toolbar", "MaterialToolbar", "@layout/toolbar_layout"], []), "NoticeBar": (["—"], []), "Card": (["CardView", "MaterialCardView"], []), "Tag": (["Chip"], ["—"]),
        "CountDown": (["—"], []), "Divider": (["1px View 细线", "—"], ["MaterialDivider"]), "Cell": (["—"], ["PreferenceRow", "ListItem"]), "ProductCard": (["CardView", "—"], []),
        "Price": (["—"], []), "ActionBar": (["—"], []), "Button": (["MaterialButton", "Button[Android SDK widgets]"], ["Button"]),
    },
    ("newpipe", "holdings"): {
        "导航栏": (["Toolbar", "MaterialToolbar", "@layout/toolbar_layout"], []), "公告通知栏": (["—"], []), "资产卡片": (["CardView", "MaterialCardView"], []), "标签页": (["TabLayout"], []),
        "持仓列表": (["RecyclerView"], []), "持仓列表项": (["—"], ["PreferenceRow", "ListItem"]), "涨跌标签": (["Chip"], ["—"]), "收盘倒计时": (["—"], []),
        "底部操作栏": (["—"], []), "买入按钮": (["MaterialButton", "Button[Android SDK widgets]"], ["Button"]), "卖出按钮": (["MaterialButton", "Button[Android SDK widgets]"], ["Button"]), "更多按钮": (["MaterialButton", "Button[Android SDK widgets]"], ["Button"]),
    },
    ("vant-demo-vue3", "order-detail"): "VANT", ("vant-demo-vue3", "holdings"): "VANT_H",
    ("vant-demo-vue2", "order-detail"): "VANT2", ("vant-demo-vue2", "holdings"): "VANT2_H",
    ("litemall", "order-detail"): "LITEMALL", ("litemall", "holdings"): "LITEMALL_H",
    ("antdm-template", "order-detail"): ANTD, ("antdm-template", "holdings"): ANTD_H,
    ("antdm-mobile", "order-detail"): ANTD, ("antdm-mobile", "holdings"): ANTD_H,
}


# third page (design/figma-export/trade-order): English layer names, not componentized; concepts the
# other pages lack (search, stepper, input, switch, checkbox, empty). Keys without "（未组件化）".
def TO(nav, search, tabs, cell, stepper, inp, switch, check, button, empty):
    return {"NavBar": nav, "SearchBar": search, "Tabs": tabs, "Cell": cell, "Stepper": stepper, "Input": inp,
            "Switch": switch, "Checkbox": check, "Button": button, "Empty": empty}
R = lambda *xs: (list(xs), [])          # right answers
RA = lambda r, a: (list(r), list(a))    # right / acceptable
NEW = R("—")
VANT_TO = TO(R("NavBar"), R("Search"), R("Tabs"), R("Cell"), R("Stepper"), R("Field"), R("Switch"), R("Checkbox"), R("Button"), R("Empty"))
EXPECTED_TO = {
    "newbee": dict(VANT_TO, NavBar=RA(["SimpleHeader"], ["NavBar"])),
    # litemall has its own is-empty (3 uses): the page check would flag a new van-empty as inconsistent
    "litemall": dict(VANT_TO, Empty=R("IsEmpty", "Empty")), "vant-demo-vue3": VANT_TO, "vant-demo-vue2": VANT_TO,
    "qmui": TO(R("QMUITopBarLayout", "QMUITopBar"), NEW, R("QMUITabSegment"), R("QMUICommonListItemView"), NEW, RA(["EditText"], ["EmojiEditText"]),
               RA(["—"], ["Switch"]), RA(["—", "CheckBox"], ["MaterialCheckBox"]), R("QMUIRoundButton"), R("QMUIEmptyView")),
    "nia": TO(R("NiaTopAppBar"), R("SearchBar", "TextField"), R("NiaTabRow"), R("ListItem"), NEW, R("TextField", "OutlinedTextField"), R("Switch"), R("Checkbox"),
              R("NiaButton"), RA(["—"], ["InterestsDetailPlaceholder"])),
    "jetsnack": TO(R("TopAppBar", "CenterAlignedTopAppBar"), R("SearchBar", "TextField"), R("TabRow"), R("ListItem"), R("QuantitySelector"), R("TextField", "OutlinedTextField"),
                   R("Switch"), R("Checkbox"), R("JetsnackButton"), NEW),
    "antdm-template": TO(R("NavBar"), R("SearchBar"), R("Tabs"), R("ListItem"), R("Stepper"), R("Input"), R("Switch"), R("Checkbox"), R("Button"), R("Empty", "ErrorBlock")),
    "hello-uniapp": TO(R("NavBar"), R("SearchBar"), R("SegmentedControl"), R("ListItem"), R("NumberBox"), R("Input", "Easyinput"), R("Switch"), R("Checkbox", "DataCheckbox"),
                       R("Button[uni-app 内置组件]"), RA(["—"], ["Empty"])),
    "taro-ui-demo": TO(R("AtNavBar"), R("AtSearchBar"), R("AtTabs"), R("AtListItem"), R("AtInputNumber"), R("AtInput"), R("AtSwitch"), R("AtCheckbox"), R("AtButton"), NEW),
    "nutui-vite": TO(R("Navbar"), R("SearchBar"), R("Tabs"), R("Cell"), R("InputNumber"), R("Input"), R("Switch"), R("Checkbox"), R("Button"), R("Empty")),
    "pokedex": TO(R("MaterialToolbar", "Toolbar"), RA(["—"], ["SearchBar"]), R("TabLayout"), NEW, NEW, R("TextInputEditText", "TextInputLayout", "EditText"),
                  R("SwitchMaterial", "SwitchCompat"), R("MaterialCheckBox", "CheckBox"), R("MaterialButton"), NEW),
    "newpipe": TO(R("Toolbar", "MaterialToolbar", "@layout/toolbar_layout"), RA(["—"], ["SearchBar"]), R("TabLayout"), RA(["—"], ["ListItem", "PreferenceRow"]), NEW,
                  R("NewPipeEditText", "EditText", "TextInputEditText"), RA(["Switch[Android SDK widgets]", "SwitchCompat"], ["Switch"]), R("CheckBox", "MaterialCheckBox"),
                  R("MaterialButton", "Button[Android SDK widgets]"), R("@layout/list_empty_view")),
}
EXPECTED_TO["antdm-mobile"] = EXPECTED_TO["antdm-template"]

# --- second held-out batch: answers written before running the current tool on them ---
# vue3-vant-mobile: Vant 4 template; its own NavBar.vue is the app-wide navbar
EXPECTED[("vue3-vant-mobile", "order-detail")] = "VANT"
EXPECTED[("vue3-vant-mobile", "holdings")] = "VANT_H"
EXPECTED_TO["vue3-vant-mobile"] = VANT_TO
# mall-app-web: uni-app; uni-ui installed but barely used; own mix-list-cell (6 uses), empty.vue, a local uni-number-box
UNI_BTN = R("Button[uni-app 内置组件]")
EXPECTED[("mall-app-web", "order-detail")] = {
    "NavBar": R("NavBar", "—"), "NoticeBar": R("NoticeBar"), "Card": R("Card", "—"), "Tag": R("Tag"), "CountDown": R("Countdown"),
    "Divider": NEW, "Cell": RA(["MixListCell"], ["ListItem"]), "ProductCard": RA(["—"], ["Card"]), "Price": NEW, "ActionBar": R("GoodsNav"), "Button": UNI_BTN}
EXPECTED[("mall-app-web", "holdings")] = {
    "导航栏": R("NavBar", "—"), "公告通知栏": R("NoticeBar"), "资产卡片": R("Card", "—"), "标签页": R("SegmentedControl"),
    "持仓列表": R("List", "—"), "持仓列表项": RA(["MixListCell"], ["ListItem"]), "涨跌标签": R("Tag"), "收盘倒计时": R("Countdown"),
    "底部操作栏": R("GoodsNav"), "买入按钮": UNI_BTN, "卖出按钮": UNI_BTN, "更多按钮": UNI_BTN}
EXPECTED_TO["mall-app-web"] = TO(R("NavBar", "—"), R("SearchBar"), R("SegmentedControl"), RA(["MixListCell"], ["ListItem"]), R("NumberBox", "UniNumberBox"),
                                 R("Input"), R("Switch"), R("Checkbox"), UNI_BTN, R("Empty"))
# WanAndroid (goweii): XML, AppCompat theme; top bar from the third-party per.goweii.actionbarex, states via
# MultiStateView, switches are SwitchCompat; no button component (buttons are styled TextViews)
WA_NAV = R("ActionBarCommon", "ActionBarEx")
WA_ROW = RA(["—"], ["@layout/rv_item_article"])
EXPECTED[("wanandroid", "order-detail")] = {
    "NavBar": WA_NAV, "NoticeBar": NEW, "Card": R("CardView"), "Tag": NEW, "CountDown": NEW, "Divider": R("1px View 细线", "—"),
    "Cell": WA_ROW, "ProductCard": R("CardView", "—"), "Price": NEW, "ActionBar": NEW, "Button": NEW}
EXPECTED[("wanandroid", "holdings")] = {
    "导航栏": WA_NAV, "公告通知栏": NEW, "资产卡片": R("CardView"), "标签页": R("TabLayout"), "持仓列表": R("RecyclerView"), "持仓列表项": WA_ROW,
    "涨跌标签": NEW, "收盘倒计时": NEW, "底部操作栏": NEW, "买入按钮": NEW, "卖出按钮": NEW, "更多按钮": NEW}
EXPECTED_TO["wanandroid"] = TO(WA_NAV, NEW, R("TabLayout"), WA_ROW, NEW, R("EditText", "InputView"), R("SwitchCompat"),
                               R("MaterialCheckBox", "CheckBox"), NEW, R("MultiStateView"))
# KunMinX Jetpack-MVVM-Best-Practice: XML + DataBinding, AppCompat theme with Material 1.9 on the
# classpath (MaterialCardView / MaterialDivider / SearchBar would crash); almost no component layer
KX_NAV = RA(["Toolbar", "—"], ["MaterialToolbar"])
KX_BTN = R("Button[Android SDK widgets]", "AppCompatButton")
EXPECTED[("kunminx", "order-detail")] = {
    "NavBar": KX_NAV, "NoticeBar": NEW, "Card": R("CardView", "—"), "Tag": NEW, "CountDown": NEW, "Divider": R("1px View 细线", "—"),
    "Cell": NEW, "ProductCard": R("CardView", "—"), "Price": NEW, "ActionBar": NEW, "Button": KX_BTN}
EXPECTED[("kunminx", "holdings")] = {
    "导航栏": KX_NAV, "公告通知栏": NEW, "资产卡片": R("CardView", "—"), "标签页": R("TabLayout"), "持仓列表": R("RecyclerView"), "持仓列表项": NEW,
    "涨跌标签": NEW, "收盘倒计时": NEW, "底部操作栏": NEW, "买入按钮": KX_BTN, "卖出按钮": KX_BTN, "更多按钮": KX_BTN}
EXPECTED_TO["kunminx"] = TO(KX_NAV, RA(["—"], ["SearchView"]), R("TabLayout"), NEW, NEW, R("EditText", "TextInputEditText", "AppCompatEditText"),
                            RA(["SwitchCompat", "Switch[Android SDK widgets]", "—"], ["SwitchMaterial"]), RA(["CheckBox", "AppCompatCheckBox", "—"], ["MaterialCheckBox"]),
                            KX_BTN, NEW)
# jeecg-uniapp: uni-app + wot-design-uni 1.13 (wd-*), a little uni-ui
EXPECTED[("jeecg-uniapp", "order-detail")] = {
    "NavBar": R("Navbar"), "NoticeBar": R("NoticeBar"), "Card": R("CellGroup", "Card"), "Tag": R("Tag"), "CountDown": R("CountDown"), "Divider": R("Divider"),
    "Cell": R("Cell"), "ProductCard": R("Card", "—"), "Price": NEW, "ActionBar": RA(["—"], ["GoodsNav"]), "Button": R("Button")}
EXPECTED[("jeecg-uniapp", "holdings")] = {
    "导航栏": R("Navbar"), "公告通知栏": R("NoticeBar"), "资产卡片": R("CellGroup", "Card"), "标签页": R("Tabs", "Segmented"), "持仓列表": R("List", "—"),
    "持仓列表项": R("Cell"), "涨跌标签": R("Tag"), "收盘倒计时": R("CountDown"), "底部操作栏": RA(["—"], ["GoodsNav"]),
    "买入按钮": R("Button"), "卖出按钮": R("Button"), "更多按钮": R("Button")}
EXPECTED_TO["jeecg-uniapp"] = TO(R("Navbar"), R("Search"), R("Tabs", "Segmented"), R("Cell"), R("InputNumber"), R("Input"), R("Switch"), R("Checkbox"),
                                 R("Button"), R("StatusTip", "Empty"))
# sunniejs/vue-h5-template: Vant / NutUI / Varlet installed, .env VITE_UI_FRAMEWORK=vant; its own shop ProductCard
EXPECTED[("vue-h5-template", "order-detail")] = dict(VANT, ProductCard=R("ProductCard", "Card"), Price=RA(["—"], ["Price"]))
EXPECTED[("vue-h5-template", "holdings")] = "VANT_H"
EXPECTED_TO["vue-h5-template"] = VANT_TO
# Seal: Compose + Material3 with its own ui/component set (SettingItem, PreferenceSwitch, SealSearchBar …)
SEAL_BTN = R("Button", "FilledButtonWithIcon", "OutlinedButtonWithIcon")
SEAL_NAV = R("TopAppBar", "LargeTopAppBar", "CenterAlignedTopAppBar")
SEAL_ROW = R("SettingItem", "PreferenceItem", "ListItem")
SEAL_TAG = RA(["—"], ["ButtonChip", "AssistChip", "SuggestionChip"])
EXPECTED[("seal", "order-detail")] = {
    "NavBar": SEAL_NAV, "NoticeBar": NEW, "Card": R("Card", "ElevatedCard", "OutlinedCard"), "Tag": SEAL_TAG, "CountDown": NEW,
    "Divider": R("HorizontalDivider", "Divider"), "Cell": SEAL_ROW, "ProductCard": RA(["Card", "—"], ["VideoCard"]), "Price": NEW, "ActionBar": NEW, "Button": SEAL_BTN}
EXPECTED[("seal", "holdings")] = {
    "导航栏": SEAL_NAV, "公告通知栏": NEW, "资产卡片": R("Card", "ElevatedCard", "OutlinedCard"), "标签页": R("TabRow", "PrimaryTabRow", "SecondaryTabRow"),
    "持仓列表": R("LazyColumn"), "持仓列表项": SEAL_ROW, "涨跌标签": SEAL_TAG, "收盘倒计时": NEW, "底部操作栏": NEW,
    "买入按钮": SEAL_BTN, "卖出按钮": SEAL_BTN, "更多按钮": SEAL_BTN}
EXPECTED_TO["seal"] = TO(SEAL_NAV, R("SealSearchBar", "SearchBar"), R("TabRow", "PrimaryTabRow", "SecondaryTabRow"), SEAL_ROW, NEW,
                         R("OutlinedTextField", "TextField"), R("Switch", "PreferenceSwitch"), R("Checkbox"), SEAL_BTN, RA(["—"], ["ErrorMessage"]))
# AntennaPod: XML + Material 1.14 with a Material3 theme; empty states go through its EmptyViewHandler (7 screens)
AP_BTN = R("MaterialButton")
EXPECTED[("antennapod", "order-detail")] = {
    "NavBar": R("MaterialToolbar", "Toolbar"), "NoticeBar": NEW, "Card": R("MaterialCardView", "CardView"), "Tag": R("Chip"), "CountDown": NEW,
    "Divider": R("MaterialDivider", "1px View 细线"), "Cell": RA(["—"], ["@layout/swipeactions_row"]), "ProductCard": R("MaterialCardView", "CardView", "—"),
    "Price": NEW, "ActionBar": NEW, "Button": AP_BTN}
EXPECTED[("antennapod", "holdings")] = {
    "导航栏": R("MaterialToolbar", "Toolbar"), "公告通知栏": NEW, "资产卡片": R("MaterialCardView", "CardView"), "标签页": R("TabLayout"),
    "持仓列表": R("RecyclerView"), "持仓列表项": RA(["—"], ["@layout/swipeactions_row"]), "涨跌标签": R("Chip"), "收盘倒计时": NEW,
    "底部操作栏": NEW, "买入按钮": AP_BTN, "卖出按钮": AP_BTN, "更多按钮": AP_BTN}
EXPECTED_TO["antennapod"] = TO(R("MaterialToolbar", "Toolbar"), R("SearchBar", "SearchView"), R("TabLayout"), RA(["—"], ["@layout/swipeactions_row"]), NEW,
                               R("TextInputEditText", "TextInputLayout", "EditText"), R("MaterialSwitch", "SwitchMaterial", "SwitchCompat"),
                               R("MaterialCheckBox", "CheckBox"), AP_BTN, R("EmptyViewHandler", "@layout/empty_view_layout"))
EXPECTED_TO["sunflower"] = TO(R("TopAppBar", "CenterAlignedTopAppBar"), R("SearchBar"), R("TabRow"), R("ListItem"), NEW, R("TextField", "OutlinedTextField"),
                              R("Switch"), R("Checkbox"), R("Button"), NEW)
EXPECTED[("sunflower", "order-detail")] = {
    "NavBar": R("TopAppBar", "CenterAlignedTopAppBar"), "NoticeBar": NEW, "Card": R("Card"), "Tag": RA(["—"], ["InputChip", "AssistChip", "SuggestionChip"]),
    "CountDown": NEW, "Divider": R("HorizontalDivider"), "Cell": R("ListItem"), "ProductCard": RA(["Card", "—"], []),
    "Price": NEW, "ActionBar": NEW, "Button": R("Button")}
EXPECTED[("sunflower", "holdings")] = {
    "导航栏": R("TopAppBar", "CenterAlignedTopAppBar"), "公告通知栏": NEW, "资产卡片": R("Card"), "标签页": R("TabRow"),
    "持仓列表": R("LazyColumn"), "持仓列表项": R("ListItem"), "涨跌标签": RA(["—"], ["InputChip", "AssistChip", "SuggestionChip"]), "收盘倒计时": NEW,
    "底部操作栏": NEW, "买入按钮": R("Button"), "卖出按钮": R("Button"), "更多按钮": R("Button")}
for _k, _v in EXPECTED_TO.items():
    EXPECTED[(_k, "trade-order")] = _v


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
        for page in ("order-detail", "holdings", "trade-order"):
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
                src = it["candidates"][0]["source"] if top != "—" else ""
                right, ok = exp.get(key, ([], []))
                hit = lambda xs: top in xs or f"{top}[{src}]" in xs  # "Button[Android SDK widgets]" pins the source
                g = "✓" if hit(right) else "△" if hit(ok) else "✗"
                counts[g] += 1
                total[g] += 1
                if g != "✓":
                    lines.append(f"    {g} {key}: {top}（期望 {' / '.join(right)}）")
            report[f"{name}/{page}"] = counts
            print(f"{name:15s} {page:13s} ✓{counts['✓']:2d} △{counts['△']} ✗{counts['✗']}")
            print("\n".join(lines)) if lines else None
    print(f"合计 ✓{total['✓']} △{total['△']} ✗{total['✗']}")
    if not only:  # a partial run must not overwrite the full results
        json.dump(report, open(os.path.join(HERE, "real-match-eval.json"), "w"), ensure_ascii=False, indent=1)
        held = {"sunflower", "pokedex", "jetsnack", "vant-demo-vue3", "vant-demo-vue2", "newpipe", "litemall", "hello-uniapp", "taro-ui-demo", "nutui-vite",
                "vue3-vant-mobile", "mall-app-web", "antennapod", "wanandroid", "seal", "kunminx", "jeecg-uniapp", "vue-h5-template"}
        lines = ["# figmafe match 回归评测（真实开源项目，固定 commit）", "",
                 "由 `python3 experiments/real-match-eval.py` 生成。✓ = 推荐正确，△ = 可用但不是最优，✗ = 推荐错误。",
                 "「留出」= 调规则时没用过、第一次跑完才修问题的项目（第一次的结果见 match-judgement.md 第 3 节）。", "",
                 "| 项目 | 技术栈 | 订单详情（组件化设计稿） | 我的持仓（未组件化、中文图层） | 买入下单（未组件化、英文图层、新概念） |", "|---|---|---|---|---|"]
        stack = {"newbee": "Vue 3 + Vant", "qmui": "XML + QMUI", "nia": "Compose + 自有设计系统", "antdm-template": "React + antd-mobile",
                 "antdm-mobile": "React + antd-mobile", "jetsnack": "Compose + 自有设计系统", "litemall": "Vue 2 + Vant 2",
                 "hello-uniapp": "uni-app + uni-ui", "taro-ui-demo": "Taro + taro-ui", "nutui-vite": "Vue 3 + NutUI",
                 "newpipe": "XML（AppCompat）+ 少量 Compose", "pokedex": "XML + Material（MaterialComponents 主题）", "sunflower": "Compose + Material3（无自有设计系统）", "vant-demo-vue3": "Vue 3 + Vant", "vant-demo-vue2": "Vue 2 + Vant 2",
                 "vue3-vant-mobile": "Vue 3 + Vant 4", "mall-app-web": "uni-app + 自有组件", "antennapod": "XML + Material3 主题",
                 "wanandroid": "XML（AppCompat）+ 三方 ActionBarEx", "seal": "Compose + Material3 + 自有组件",
                 "kunminx": "XML + DataBinding（AppCompat 主题）",
                 "jeecg-uniapp": "uni-app + wot-design-uni",
                 "vue-h5-template": "Vue 3 + Vant（另装 NutUI / Varlet）"}
        cell = lambda c: f"✓{c['✓']} △{c['△']} ✗{c['✗']}" if c else "—"
        for name in PROJECTS:
            tag = "（留出）" if name in held else ""
            lines.append(f"| {name}{tag} | {stack.get(name, '')} | {cell(report.get(name + '/order-detail'))} | {cell(report.get(name + '/holdings'))} | {cell(report.get(name + '/trade-order'))} |")
        lines += ["", f"**合计 ✓{total['✓']} / △{total['△']} / ✗{total['✗']}**", ""]
        open(os.path.join(HERE, "real-match-eval.md"), "w").write("\n".join(lines))
    sys.exit(1 if total["✗"] else 0)


if __name__ == "__main__":
    main()
