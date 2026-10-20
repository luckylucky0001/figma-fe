package reuse

// Built-in catalogs of common mobile UI libraries, so a project that depends on one of them gets
// those components as reusable assets even when node_modules / AARs are not on disk.
// Names are the public component names; usage is detected per-library (tag prefix or import).

type uiLib struct {
	Name       string   // display name
	Packages   []string // npm package or Gradle group to detect
	TagPrefix  string   // Vue kebab-case tag prefix (van-, nut-, t-) or "" for React/Android
	Import     []string // import sources (React / Vue named imports)
	Components []string // PascalCase component names
	Platform   string
}

var uiLibs = []uiLib{
	{Name: "Vant", Packages: []string{"vant"}, TagPrefix: "van-", Import: []string{"vant"}, Platform: "h5", Components: []string{
		"ActionBar", "ActionBarButton", "ActionBarIcon", "GoodsAction", "GoodsActionButton", "GoodsActionIcon", "ActionSheet", "AddressEdit", "AddressList", "Area", "BackTop", "Badge", "Button", "Calendar", "Card", "Cascader", "Cell", "CellGroup", "Checkbox", "CheckboxGroup", "Circle", "Col", "Collapse", "CollapseItem", "ConfigProvider", "ContactCard", "ContactEdit", "ContactList", "CountDown", "Coupon", "CouponCell", "CouponList", "DatePicker", "Dialog", "Divider", "DropdownItem", "DropdownMenu", "Empty", "Field", "FloatingBubble", "FloatingPanel", "Form", "Grid", "GridItem", "Highlight", "Icon", "Image", "ImagePreview", "IndexAnchor", "IndexBar", "List", "Loading", "NavBar", "NoticeBar", "Notify", "NumberKeyboard", "Overlay", "Pagination", "PasswordInput", "Picker", "PickerGroup", "Popover", "Popup", "Progress", "PullRefresh", "Radio", "RadioGroup", "Rate", "RollingText", "Row", "Search", "ShareSheet", "Sidebar", "SidebarItem", "Signature", "Skeleton", "Slider", "Space", "Step", "Stepper", "Steps", "Sticky", "SubmitBar", "Swipe", "SwipeCell", "SwipeItem", "Switch", "Tab", "Tabbar", "TabbarItem", "Tabs", "Tag", "TextEllipsis", "TimePicker", "Toast", "TreeSelect", "Uploader", "Watermark"}},
	{Name: "React Vant", Packages: []string{"react-vant"}, Import: []string{"react-vant"}, Platform: "h5", Components: []string{
		"ActionBar", "ActionSheet", "Area", "Badge", "Button", "Calendar", "Card", "Cascader", "Cell", "Checkbox", "Circle", "Collapse", "ConfigProvider", "CountDown", "DatetimePicker", "Dialog", "Divider", "DropdownMenu", "Empty", "Field", "Flex", "Form", "Grid", "Hooks", "Image", "ImagePreview", "IndexBar", "Input", "List", "Loading", "NavBar", "NoticeBar", "NumberKeyboard", "Overlay", "PasswordInput", "Picker", "Popover", "Popup", "Progress", "PullRefresh", "Radio", "Rate", "Search", "ShareSheet", "Sidebar", "Skeleton", "Slider", "Space", "Steps", "Sticky", "Stepper", "SubmitBar", "Swiper", "SwipeCell", "Switch", "Tabbar", "Tabs", "Tag", "TextEllipsis", "Toast", "TreeSelect", "Typography", "Uploader"}},
	{Name: "NutUI", Packages: []string{"@nutui/nutui", "@nutui/nutui-react", "@nutui/nutui-taro"}, TagPrefix: "nut-", Import: []string{"@nutui/nutui", "@nutui/nutui-react"}, Platform: "h5", Components: []string{
		"ActionSheet", "Address", "Animate", "Avatar", "Backtop", "Badge", "Barrage", "Button", "Calendar", "Card", "Cascader", "Cell", "CellGroup", "Checkbox", "CircleProgress", "Collapse", "Countdown", "CountUp", "DatePicker", "Dialog", "Divider", "Drag", "Ecard", "Elevator", "Empty", "FixedNav", "Form", "Grid", "Icon", "Image", "ImagePreview", "Indicator", "InfiniteLoading", "Input", "InputNumber", "Menu", "Navbar", "NoticeBar", "Notify", "NumberKeyboard", "Overlay", "Pagination", "Picker", "Popover", "Popup", "Price", "Progress", "Radio", "Range", "Rate", "SearchBar", "ShortPassword", "SideNavBar", "Signature", "Skeleton", "Sku", "Steps", "Sticky", "Swipe", "Swiper", "Switch", "Tabbar", "Tabs", "Tag", "TextArea", "Toast", "TrendArrow", "Uploader", "Watermark"}},
	{Name: "antd-mobile", Packages: []string{"antd-mobile"}, Import: []string{"antd-mobile"}, Platform: "h5", Components: []string{
		"InputItem", "TextareaItem", "Carousel", "ActivityIndicator", "WhiteSpace", "WingBlank", "ListView", "Flex", "Accordion", "Range", "Menu", "ListItem", "FormItem", "TabBarItem", "SwiperItem", "GridItem", "CollapsePanel", "ActionSheet", "AutoCenter", "Avatar", "Badge", "Button", "Calendar", "CapsuleTabs", "Card", "CascadePicker", "Cascader", "CenterPopup", "CheckList", "Checkbox", "Collapse", "DatePicker", "Dialog", "Divider", "DotLoading", "Dropdown", "Ellipsis", "Empty", "ErrorBlock", "FloatingBubble", "FloatingPanel", "Form", "Grid", "Image", "ImageUploader", "ImageViewer", "IndexBar", "InfiniteScroll", "Input", "JumboTabs", "List", "Loading", "Mask", "Modal", "NavBar", "NoticeBar", "NumberKeyboard", "PageIndicator", "PasscodeInput", "Picker", "Popover", "Popup", "ProgressBar", "ProgressCircle", "PullToRefresh", "Radio", "Rate", "Result", "SafeArea", "ScrollMask", "SearchBar", "Selector", "SideBar", "Skeleton", "Slider", "Space", "SpinLoading", "Stepper", "Steps", "SwipeAction", "Swiper", "Switch", "TabBar", "Tabs", "Tag", "TextArea", "Toast", "TreeSelect", "VirtualInput", "WaterMark"}},
	{Name: "TDesign Mobile", Packages: []string{"tdesign-mobile-vue", "tdesign-mobile-react"}, TagPrefix: "t-", Import: []string{"tdesign-mobile-vue", "tdesign-mobile-react"}, Platform: "h5", Components: []string{
		"ActionSheet", "Avatar", "BackTop", "Badge", "Button", "Calendar", "Cascader", "Cell", "CellGroup", "Checkbox", "Collapse", "CountDown", "DateTimePicker", "Dialog", "Divider", "Drawer", "DropdownMenu", "Empty", "Fab", "Form", "Grid", "Icon", "Image", "ImageViewer", "Indexes", "Input", "Link", "Loading", "Message", "Navbar", "NoticeBar", "Overlay", "Picker", "Popover", "Popup", "Progress", "PullDownRefresh", "Radio", "Rate", "Result", "Search", "SideBar", "Skeleton", "Slider", "Stepper", "Steps", "Sticky", "SwipeCell", "Swiper", "Switch", "TabBar", "Tabs", "Tag", "Textarea", "Toast", "Upload"}},
	// uni-app's built-in components are lowercase tags (<button>, <navigator>); active when the project has pages.json
	{Name: "uni-app 内置组件", Packages: []string{"@dcloudio/uni-app"}, Platform: "h5", Components: []string{
		"Button", "Navigator", "ScrollView", "Swiper", "Switch", "Checkbox", "Radio", "Picker", "Progress", "Image", "Input", "Textarea", "Slider", "Icon"}},
	{Name: "taro-ui", Packages: []string{"taro-ui", "taro-ui-vue3"}, Import: []string{"taro-ui", "taro-ui-vue3"}, Platform: "h5", Components: []string{
		"AtActionSheet", "AtActionSheetItem", "AtActivityIndicator", "AtAccordion", "AtAvatar", "AtBadge", "AtButton", "AtCalendar", "AtCard", "AtCheckbox",
		"AtCountdown", "AtCurtain", "AtDivider", "AtDrawer", "AtFab", "AtFloatLayout", "AtForm", "AtGrid", "AtIcon", "AtImagePicker", "AtIndexes", "AtInput",
		"AtInputNumber", "AtList", "AtListItem", "AtLoadMore", "AtMessage", "AtModal", "AtModalAction", "AtModalContent", "AtModalHeader", "AtNavBar",
		"AtNoticebar", "AtPagination", "AtProgress", "AtRadio", "AtRange", "AtRate", "AtSearchBar", "AtSegmentedControl", "AtSlider", "AtSteps",
		"AtSwipeAction", "AtSwitch", "AtTabBar", "AtTabs", "AtTabsPane", "AtTag", "AtTextarea", "AtTimeline", "AtToast"}},
	{Name: "uni-ui", Packages: []string{"@dcloudio/uni-ui"}, TagPrefix: "uni-", Import: []string{"@dcloudio/uni-ui"}, Platform: "h5", Components: []string{
		"Badge", "Breadcrumb", "Calendar", "Card", "Collapse", "CollapseItem", "Combox", "Countdown", "DataCheckbox", "DataPicker", "DatetimePicker", "Drawer",
		"Easyinput", "Fab", "Fav", "FilePicker", "Forms", "FormsItem", "GoodsNav", "Grid", "GridItem", "Group", "Icons", "IndexedList", "Link", "List",
		"ListItem", "ListChat", "LoadMore", "NavBar", "NoticeBar", "NumberBox", "Pagination", "Popup", "PopupDialog", "PopupMessage", "PopupShare", "Rate",
		"Row", "Col", "SearchBar", "SegmentedControl", "Section", "Steps", "SwipeAction", "SwipeActionItem", "SwiperDot", "Table", "Tag", "Title", "Tooltip", "Transition"}},
	{Name: "wot-design-uni", Packages: []string{"wot-design-uni"}, TagPrefix: "wd-", Import: []string{"wot-design-uni"}, Platform: "h5", Components: []string{
		"ActionSheet", "Avatar", "AvatarGroup", "Backtop", "Badge", "Button", "Calendar", "CalendarView", "Card", "Cell", "CellGroup", "Checkbox", "CheckboxGroup", "Circle",
		"Col", "ColPicker", "Collapse", "CollapseItem", "ConfigProvider", "CountDown", "CountTo", "Curtain", "DatetimePicker", "DatetimePickerView", "Divider", "DropMenu", "DropMenuItem",
		"Fab", "FloatingPanel", "Form", "FormItem", "Gap", "Grid", "GridItem", "Icon", "Img", "ImgCropper", "IndexAnchor", "IndexBar", "Input", "InputNumber", "Keyboard", "Loading",
		"Loadmore", "MessageBox", "Navbar", "NavbarCapsule", "NoticeBar", "Notify", "NumberKeyboard", "Overlay", "Pagination", "PasswordInput", "Picker", "PickerView", "Popover", "Popup",
		"Progress", "Radio", "RadioGroup", "Rate", "Resize", "Row", "Search", "Segmented", "SelectPicker", "Sidebar", "SidebarItem", "Signature", "Skeleton", "Slider", "SortButton",
		"StatusTip", "Step", "Steps", "Sticky", "StickyBox", "SwipeAction", "Swiper", "SwiperNav", "Switch", "Tab", "Tabbar", "TabbarItem", "Table", "TableCol", "Tabs", "Tag", "Text",
		"Textarea", "Toast", "Tooltip", "Tour", "Transition", "Upload", "Watermark"}},
	{Name: "uView", Packages: []string{"uview-ui", "uview-plus"}, TagPrefix: "u-", Import: []string{"uview-ui", "uview-plus"}, Platform: "h5", Components: []string{
		"ActionSheet", "Album", "Alert", "Avatar", "BackTop", "Badge", "Button", "Calendar", "Card", "Cell", "CellGroup", "Checkbox", "CodeInput", "Col", "Collapse", "CountDown", "CountTo", "DatetimePicker", "Divider", "Empty", "Form", "FormItem", "Gap", "Grid", "GridItem", "Icon", "Image", "IndexList", "Input", "Keyboard", "Line", "LineProgress", "Link", "List", "ListItem", "Loadmore", "LoadingIcon", "LoadingPage", "Modal", "Navbar", "NoNetwork", "NoticeBar", "Notify", "NumberBox", "Overlay", "Parse", "Picker", "Popup", "Radio", "Rate", "ReadMore", "Row", "RowNotice", "Search", "Skeleton", "Slider", "Status", "Steps", "Sticky", "Subsection", "SwipeAction", "Swiper", "Switch", "Tabbar", "Tabs", "Tag", "Text", "Textarea", "Toast", "Tooltip", "Transition", "Upload"}},
	{Name: "Element Plus", Packages: []string{"element-plus", "element-ui"}, TagPrefix: "el-", Import: []string{"element-plus", "element-ui"}, Platform: "h5", Components: []string{
		"Alert", "Avatar", "Backtop", "Badge", "Breadcrumb", "Button", "Calendar", "Card", "Carousel", "Cascader", "Checkbox", "Col", "Collapse", "ColorPicker", "DatePicker", "Descriptions", "Dialog", "Divider", "Drawer", "Dropdown", "Empty", "Form", "FormItem", "Icon", "Image", "Input", "InputNumber", "Link", "Loading", "Menu", "Message", "MessageBox", "Notification", "Pagination", "Popconfirm", "Popover", "Progress", "Radio", "Rate", "Result", "Row", "Select", "Skeleton", "Slider", "Space", "Statistic", "Steps", "Switch", "Table", "Tabs", "Tag", "TimePicker", "Timeline", "Tooltip", "Transfer", "Tree", "Upload"}},
	{Name: "Material Components", Packages: []string{"com.google.android.material"}, Import: []string{"com.google.android.material."}, Platform: "android", Components: []string{
		"MaterialButton", "MaterialCardView", "Chip", "ChipGroup", "MaterialToolbar", "AppBarLayout", "CollapsingToolbarLayout", "BottomNavigationView", "NavigationView", "TabLayout", "TextInputLayout", "TextInputEditText", "FloatingActionButton", "ExtendedFloatingActionButton", "Snackbar", "BottomSheetDialog", "MaterialDivider", "MaterialSwitch", "SwitchMaterial", "MaterialCheckBox", "MaterialRadioButton", "Slider", "CircularProgressIndicator", "LinearProgressIndicator", "BadgeDrawable", "ShapeableImageView", "MaterialAlertDialogBuilder", "MaterialButtonToggleGroup", "SearchBar", "SearchView"}},
	{Name: "QMUI Android", Packages: []string{"com.qmuiteam"}, Import: []string{"com.qmuiteam.qmui."}, Platform: "android", Components: []string{
		"QMUITopBar", "QMUITopBarLayout", "QMUIRoundButton", "QMUIAlphaButton", "QMUICommonListItemView", "QMUIGroupListView", "QMUIEmptyView", "QMUITipDialog", "QMUIDialog", "QMUIBottomSheet", "QMUITabSegment", "QMUIRadiusImageView", "QMUIProgressBar", "QMUILoadingView", "QMUIPopup", "QMUIFloatLayout", "QMUISpanTouchFixTextView", "QMUIQQFaceView", "QMUIPullRefreshLayout", "QMUISlider", "QMUIRoundFrameLayout", "QMUIRoundLinearLayout", "QMUIContinuousNestedScrollLayout", "QMUISkinManager", "QMUISeekBar", "QMUINotchHelper"}},
	{Name: "Compose Foundation", Packages: []string{"androidx.compose.foundation.lazy", "androidx.compose.foundation.pager"}, Import: []string{"androidx.compose.foundation.lazy."}, Platform: "android", Components: []string{
		"LazyColumn", "LazyRow", "LazyVerticalGrid", "LazyHorizontalGrid", "LazyVerticalStaggeredGrid", "HorizontalPager", "VerticalPager"}},
	{Name: "Compose Material3", Packages: []string{"androidx.compose.material3"}, Import: []string{"androidx.compose.material3."}, Platform: "android", Components: []string{
		"Button", "OutlinedButton", "TextButton", "FilledTonalButton", "ElevatedButton", "Card", "OutlinedCard", "ElevatedCard", "TopAppBar", "CenterAlignedTopAppBar", "BottomAppBar", "NavigationBar", "AssistChip", "FilterChip", "SuggestionChip", "InputChip", "Badge", "Divider", "HorizontalDivider", "ListItem", "TextField", "OutlinedTextField", "Switch", "Checkbox", "RadioButton", "Slider", "AlertDialog", "ModalBottomSheet", "Snackbar", "CircularProgressIndicator", "LinearProgressIndicator", "TabRow", "PrimaryTabRow", "SecondaryTabRow", "ScrollableTabRow", "PrimaryScrollableTabRow", "SecondaryScrollableTabRow", "Tab", "FloatingActionButton", "SearchBar"}},
	{Name: "Android SDK widgets", Packages: []string{"android"}, Platform: "android", Components: []string{
		"Button", "ImageButton", "TextView", "EditText", "ImageView", "CheckBox", "RadioButton", "Switch", "ProgressBar", "SeekBar", "Spinner", "RecyclerView", "ListView", "ScrollView", "Toolbar", "CardView", "ViewPager", "ViewPager2"}},
}

// rawH5 are native elements that have a component equivalent in any UI library.
var rawH5 = map[string]string{"button": "button", "input": "input", "select": "input", "textarea": "input", "img": "image"}

// rawAndroid are framework widgets that projects usually wrap (only flagged when the project has a component of that concept).
var rawAndroid = map[string]string{"Button": "button", "ImageButton": "button", "EditText": "input", "ImageView": "image", "CheckBox": "checkbox", "RadioButton": "radio", "Switch": "switch", "ProgressBar": "loading", "Toolbar": "navbar"}

// libConcept overrides the name-derived concept where a library uses a misleading name.
var libConcept = map[string]string{
	"Vant/Card":                                    "goods", // van-card is a product card
	"antd-mobile/InputItem":                        "input",
	"antd-mobile/TextareaItem":                     "input",
	"antd-mobile/Carousel":                         "carousel",
	"antd-mobile/ActivityIndicator":                "loading",
	"antd-mobile/WhiteSpace":                       "",
	"antd-mobile/WingBlank":                        "",
	"antd-mobile/Flex":                             "",
	"antd-mobile/ListItem":                         "cell",
	"antd-mobile/FormItem":                         "", // a form row (label + control), not a list cell
	"antd-mobile/GridItem":                         "",
	"antd-mobile/CollapsePanel":                    "",
	"antd-mobile/ErrorBlock":                       "empty",
	"antd-mobile/Result":                           "empty",
	"Vant/GoodsAction":                             "actionbar",
	"Vant/GoodsActionButton":                       "button",
	"Vant/GoodsActionIcon":                         "",
	"uni-app 内置组件/Navigator":                       "",
	"uni-app 内置组件/ScrollView":                      "",
	"uni-app 内置组件/Picker":                          "",
	"uni-app 内置组件/Slider":                          "",
	"taro-ui/AtListItem":                           "cell",
	"taro-ui/AtList":                               "card",
	"taro-ui/AtNoticebar":                          "notice",
	"taro-ui/AtSegmentedControl":                   "tabs",
	"taro-ui/AtActivityIndicator":                  "loading",
	"taro-ui/AtLoadMore":                           "loading",
	"taro-ui/AtTabsPane":                           "",
	"taro-ui/AtFloatLayout":                        "popup",
	"taro-ui/AtInputNumber":                        "stepper",
	"NutUI/Sku":                                    "", // a spec picker popup
	"NutUI/Ecard":                                  "",
	"NutUI/InputNumber":                            "stepper",
	"wot-design-uni/StatusTip":                     "empty",
	"wot-design-uni/InputNumber":                   "stepper",
	"wot-design-uni/Img":                           "image",
	"wot-design-uni/Segmented":                     "tabs",
	"wot-design-uni/Loadmore":                      "loading",
	"wot-design-uni/CellGroup":                     "card",
	"wot-design-uni/Text":                          "",
	"wot-design-uni/TabbarItem":                    "",
	"wot-design-uni/PasswordInput":                 "",
	"wot-design-uni/Row":                           "",
	"wot-design-uni/Col":                           "",
	"wot-design-uni/GridItem":                      "",
	"wot-design-uni/DropMenuItem":                  "",
	"uni-ui/ListItem":                              "cell",
	"uni-ui/Section":                               "",
	"uni-ui/Title":                                 "",
	"uni-ui/Group":                                 "card",
	"uni-ui/GoodsNav":                              "actionbar",
	"uni-ui/SegmentedControl":                      "tabs",
	"uni-ui/Easyinput":                             "input",
	"uni-ui/NumberBox":                             "stepper",
	"uni-ui/LoadMore":                              "loading",
	"uni-ui/SwiperDot":                             "carousel",
	"uni-ui/Row":                                   "",
	"uni-ui/Col":                                   "",
	"Vant/SubmitBar":                               "actionbar", // price + submit: a specialised action bar
	"Vant/Row":                                     "",          // grid layout, not a list row
	"Material Components/BadgeDrawable":            "",          // a drawable overlay, not a View to place
	"Vant/Col":                                     "",
	"Vant/Space":                                   "",
	"React Vant/Card":                              "goods",
	"Vant/ActionBar":                               "actionbar",
	"Vant/Field":                                   "input",
	"Vant/CellGroup":                               "card",
	"uView/Row":                                    "",
	"uView/Col":                                    "",
	"uView/Line":                                   "divider",
	"uView/Navbar":                                 "navbar",
	"Element Plus/Row":                             "",
	"Element Plus/Col":                             "",
	"Element Plus/Descriptions":                    "cell",
	"NutUI/Card":                                   "goods",
	"NutUI/Price":                                  "price",
	"NutUI/Countdown":                              "countdown",
	"TDesign Mobile/Navbar":                        "navbar",
	"QMUI Android/QMUICommonListItemView":          "cell",
	"QMUI Android/QMUIGroupListView":               "card",
	"QMUI Android/QMUIRoundButton":                 "button",
	"QMUI Android/QMUITopBar":                      "navbar",
	"QMUI Android/QMUITopBarLayout":                "navbar",
	"Material Components/MaterialToolbar":          "navbar",
	"Material Components/Chip":                     "tag",
	"Compose Foundation/LazyRow":                   "list",
	"Compose Foundation/LazyVerticalGrid":          "",
	"Compose Foundation/LazyHorizontalGrid":        "",
	"Compose Foundation/LazyVerticalStaggeredGrid": "",
	"Compose Foundation/HorizontalPager":           "carousel",
	"Compose Foundation/VerticalPager":             "carousel",
	"Compose Material3/ListItem":                   "cell",
	"Compose Material3/AssistChip":                 "tag",
}
