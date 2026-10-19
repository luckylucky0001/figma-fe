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
		"ActionBar", "ActionBarButton", "ActionBarIcon", "ActionSheet", "AddressEdit", "AddressList", "Area", "BackTop", "Badge", "Button", "Calendar", "Card", "Cascader", "Cell", "CellGroup", "Checkbox", "CheckboxGroup", "Circle", "Col", "Collapse", "CollapseItem", "ConfigProvider", "ContactCard", "ContactEdit", "ContactList", "CountDown", "Coupon", "CouponCell", "CouponList", "DatePicker", "Dialog", "Divider", "DropdownItem", "DropdownMenu", "Empty", "Field", "FloatingBubble", "FloatingPanel", "Form", "Grid", "GridItem", "Highlight", "Icon", "Image", "ImagePreview", "IndexAnchor", "IndexBar", "List", "Loading", "NavBar", "NoticeBar", "Notify", "NumberKeyboard", "Overlay", "Pagination", "PasswordInput", "Picker", "PickerGroup", "Popover", "Popup", "Progress", "PullRefresh", "Radio", "RadioGroup", "Rate", "RollingText", "Row", "Search", "ShareSheet", "Sidebar", "SidebarItem", "Signature", "Skeleton", "Slider", "Space", "Step", "Stepper", "Steps", "Sticky", "SubmitBar", "Swipe", "SwipeCell", "SwipeItem", "Switch", "Tab", "Tabbar", "TabbarItem", "Tabs", "Tag", "TextEllipsis", "TimePicker", "Toast", "TreeSelect", "Uploader", "Watermark"}},
	{Name: "React Vant", Packages: []string{"react-vant"}, Import: []string{"react-vant"}, Platform: "h5", Components: []string{
		"ActionBar", "ActionSheet", "Area", "Badge", "Button", "Calendar", "Card", "Cascader", "Cell", "Checkbox", "Circle", "Collapse", "ConfigProvider", "CountDown", "DatetimePicker", "Dialog", "Divider", "DropdownMenu", "Empty", "Field", "Flex", "Form", "Grid", "Hooks", "Image", "ImagePreview", "IndexBar", "Input", "List", "Loading", "NavBar", "NoticeBar", "NumberKeyboard", "Overlay", "PasswordInput", "Picker", "Popover", "Popup", "Progress", "PullRefresh", "Radio", "Rate", "Search", "ShareSheet", "Sidebar", "Skeleton", "Slider", "Space", "Steps", "Sticky", "Stepper", "SubmitBar", "Swiper", "SwipeCell", "Switch", "Tabbar", "Tabs", "Tag", "TextEllipsis", "Toast", "TreeSelect", "Typography", "Uploader"}},
	{Name: "NutUI", Packages: []string{"@nutui/nutui", "@nutui/nutui-react", "@nutui/nutui-taro"}, TagPrefix: "nut-", Import: []string{"@nutui/nutui", "@nutui/nutui-react"}, Platform: "h5", Components: []string{
		"ActionSheet", "Address", "Animate", "Avatar", "Backtop", "Badge", "Barrage", "Button", "Calendar", "Card", "Cascader", "Cell", "CellGroup", "Checkbox", "CircleProgress", "Collapse", "Countdown", "CountUp", "DatePicker", "Dialog", "Divider", "Drag", "Ecard", "Elevator", "Empty", "FixedNav", "Form", "Grid", "Icon", "Image", "ImagePreview", "Indicator", "InfiniteLoading", "Input", "InputNumber", "Menu", "Navbar", "NoticeBar", "Notify", "NumberKeyboard", "Overlay", "Pagination", "Picker", "Popover", "Popup", "Price", "Progress", "Radio", "Range", "Rate", "SearchBar", "ShortPassword", "SideNavBar", "Signature", "Skeleton", "Sku", "Steps", "Sticky", "Swipe", "Swiper", "Switch", "Tabbar", "Tabs", "Tag", "TextArea", "Toast", "TrendArrow", "Uploader", "Watermark"}},
	{Name: "antd-mobile", Packages: []string{"antd-mobile"}, Import: []string{"antd-mobile"}, Platform: "h5", Components: []string{
		"ActionSheet", "AutoCenter", "Avatar", "Badge", "Button", "Calendar", "CapsuleTabs", "Card", "CascadePicker", "Cascader", "CenterPopup", "CheckList", "Checkbox", "Collapse", "DatePicker", "Dialog", "Divider", "DotLoading", "Dropdown", "Ellipsis", "Empty", "ErrorBlock", "FloatingBubble", "FloatingPanel", "Form", "Grid", "Image", "ImageUploader", "ImageViewer", "IndexBar", "InfiniteScroll", "Input", "JumboTabs", "List", "Loading", "Mask", "Modal", "NavBar", "NoticeBar", "NumberKeyboard", "PageIndicator", "PasscodeInput", "Picker", "Popover", "Popup", "ProgressBar", "ProgressCircle", "PullToRefresh", "Radio", "Rate", "Result", "SafeArea", "ScrollMask", "SearchBar", "Selector", "SideBar", "Skeleton", "Slider", "Space", "SpinLoading", "Stepper", "Steps", "SwipeAction", "Swiper", "Switch", "TabBar", "Tabs", "Tag", "TextArea", "Toast", "TreeSelect", "VirtualInput", "WaterMark"}},
	{Name: "TDesign Mobile", Packages: []string{"tdesign-mobile-vue", "tdesign-mobile-react"}, TagPrefix: "t-", Import: []string{"tdesign-mobile-vue", "tdesign-mobile-react"}, Platform: "h5", Components: []string{
		"ActionSheet", "Avatar", "BackTop", "Badge", "Button", "Calendar", "Cascader", "Cell", "CellGroup", "Checkbox", "Collapse", "CountDown", "DateTimePicker", "Dialog", "Divider", "Drawer", "DropdownMenu", "Empty", "Fab", "Form", "Grid", "Icon", "Image", "ImageViewer", "Indexes", "Input", "Link", "Loading", "Message", "Navbar", "NoticeBar", "Overlay", "Picker", "Popover", "Popup", "Progress", "PullDownRefresh", "Radio", "Rate", "Result", "Search", "SideBar", "Skeleton", "Slider", "Stepper", "Steps", "Sticky", "SwipeCell", "Swiper", "Switch", "TabBar", "Tabs", "Tag", "Textarea", "Toast", "Upload"}},
	{Name: "uView", Packages: []string{"uview-ui", "uview-plus"}, TagPrefix: "u-", Import: []string{"uview-ui", "uview-plus"}, Platform: "h5", Components: []string{
		"ActionSheet", "Album", "Alert", "Avatar", "BackTop", "Badge", "Button", "Calendar", "Card", "Cell", "CellGroup", "Checkbox", "CodeInput", "Col", "Collapse", "CountDown", "CountTo", "DatetimePicker", "Divider", "Empty", "Form", "FormItem", "Gap", "Grid", "GridItem", "Icon", "Image", "IndexList", "Input", "Keyboard", "Line", "LineProgress", "Link", "List", "ListItem", "Loadmore", "LoadingIcon", "LoadingPage", "Modal", "Navbar", "NoNetwork", "NoticeBar", "Notify", "NumberBox", "Overlay", "Parse", "Picker", "Popup", "Radio", "Rate", "ReadMore", "Row", "RowNotice", "Search", "Skeleton", "Slider", "Status", "Steps", "Sticky", "Subsection", "SwipeAction", "Swiper", "Switch", "Tabbar", "Tabs", "Tag", "Text", "Textarea", "Toast", "Tooltip", "Transition", "Upload"}},
	{Name: "Element Plus", Packages: []string{"element-plus", "element-ui"}, TagPrefix: "el-", Import: []string{"element-plus", "element-ui"}, Platform: "h5", Components: []string{
		"Alert", "Avatar", "Backtop", "Badge", "Breadcrumb", "Button", "Calendar", "Card", "Carousel", "Cascader", "Checkbox", "Col", "Collapse", "ColorPicker", "DatePicker", "Descriptions", "Dialog", "Divider", "Drawer", "Dropdown", "Empty", "Form", "FormItem", "Icon", "Image", "Input", "InputNumber", "Link", "Loading", "Menu", "Message", "MessageBox", "Notification", "Pagination", "Popconfirm", "Popover", "Progress", "Radio", "Rate", "Result", "Row", "Select", "Skeleton", "Slider", "Space", "Statistic", "Steps", "Switch", "Table", "Tabs", "Tag", "TimePicker", "Timeline", "Tooltip", "Transfer", "Tree", "Upload"}},
	{Name: "Material Components", Packages: []string{"com.google.android.material"}, Import: []string{"com.google.android.material."}, Platform: "android", Components: []string{
		"MaterialButton", "MaterialCardView", "Chip", "ChipGroup", "MaterialToolbar", "AppBarLayout", "CollapsingToolbarLayout", "BottomNavigationView", "NavigationView", "TabLayout", "TextInputLayout", "TextInputEditText", "FloatingActionButton", "ExtendedFloatingActionButton", "Snackbar", "BottomSheetDialog", "MaterialDivider", "MaterialSwitch", "SwitchMaterial", "MaterialCheckBox", "MaterialRadioButton", "Slider", "CircularProgressIndicator", "LinearProgressIndicator", "BadgeDrawable", "ShapeableImageView", "MaterialAlertDialogBuilder", "MaterialButtonToggleGroup"}},
	{Name: "QMUI Android", Packages: []string{"com.qmuiteam"}, Import: []string{"com.qmuiteam.qmui."}, Platform: "android", Components: []string{
		"QMUITopBar", "QMUITopBarLayout", "QMUIRoundButton", "QMUIAlphaButton", "QMUICommonListItemView", "QMUIGroupListView", "QMUIEmptyView", "QMUITipDialog", "QMUIDialog", "QMUIBottomSheet", "QMUITabSegment", "QMUIRadiusImageView", "QMUIProgressBar", "QMUILoadingView", "QMUIPopup", "QMUIFloatLayout", "QMUISpanTouchFixTextView", "QMUIQQFaceView", "QMUIPullRefreshLayout", "QMUISlider", "QMUIRoundFrameLayout", "QMUIRoundLinearLayout", "QMUIContinuousNestedScrollLayout", "QMUISkinManager", "QMUISeekBar", "QMUINotchHelper"}},
	{Name: "Compose Material3", Packages: []string{"androidx.compose.material3"}, Import: []string{"androidx.compose.material3."}, Platform: "android", Components: []string{
		"Button", "OutlinedButton", "TextButton", "FilledTonalButton", "ElevatedButton", "Card", "OutlinedCard", "ElevatedCard", "TopAppBar", "CenterAlignedTopAppBar", "BottomAppBar", "NavigationBar", "AssistChip", "FilterChip", "SuggestionChip", "InputChip", "Badge", "Divider", "HorizontalDivider", "ListItem", "TextField", "OutlinedTextField", "Switch", "Checkbox", "RadioButton", "Slider", "AlertDialog", "ModalBottomSheet", "Snackbar", "CircularProgressIndicator", "LinearProgressIndicator", "TabRow", "Tab", "FloatingActionButton", "SearchBar"}},
	{Name: "Android SDK widgets", Packages: []string{"android"}, Platform: "android", Components: []string{
		"Button", "ImageButton", "TextView", "EditText", "ImageView", "CheckBox", "RadioButton", "Switch", "ProgressBar", "SeekBar", "Spinner", "RecyclerView", "ListView", "ScrollView", "Toolbar", "CardView", "ViewPager", "ViewPager2"}},
}

// rawH5 are native elements that have a component equivalent in any UI library.
var rawH5 = map[string]string{"button": "button", "input": "input", "select": "input", "textarea": "input", "img": "image"}

// rawAndroid are framework widgets that projects usually wrap (only flagged when the project has a component of that concept).
var rawAndroid = map[string]string{"Button": "button", "ImageButton": "button", "EditText": "input", "ImageView": "image", "CheckBox": "switch", "Switch": "switch", "ProgressBar": "loading", "Toolbar": "navbar"}

// libConcept overrides the name-derived concept where a library uses a misleading name.
var libConcept = map[string]string{
	"Vant/Card":                           "goods", // van-card is a product card
	"Vant/Row":                            "",      // grid layout, not a list row
	"Vant/Col":                            "",
	"Vant/Space":                          "",
	"React Vant/Card":                     "goods",
	"Vant/SubmitBar":                      "actionbar",
	"Vant/ActionBar":                      "actionbar",
	"Vant/Field":                          "input",
	"Vant/CellGroup":                      "card",
	"uView/Row":                           "",
	"uView/Col":                           "",
	"uView/Line":                          "divider",
	"uView/Navbar":                        "navbar",
	"Element Plus/Row":                    "",
	"Element Plus/Col":                    "",
	"Element Plus/Descriptions":           "cell",
	"NutUI/Card":                          "goods",
	"NutUI/Price":                         "price",
	"NutUI/Countdown":                     "countdown",
	"TDesign Mobile/Navbar":               "navbar",
	"QMUI Android/QMUICommonListItemView": "cell",
	"QMUI Android/QMUIGroupListView":      "card",
	"QMUI Android/QMUIRoundButton":        "button",
	"QMUI Android/QMUITopBar":             "navbar",
	"QMUI Android/QMUITopBarLayout":       "navbar",
	"Material Components/MaterialToolbar": "navbar",
	"Material Components/Chip":            "tag",
	"Compose Material3/ListItem":          "cell",
	"Compose Material3/AssistChip":        "tag",
}
