# ds-view — legacy XML View components (source-only)

Source-only legacy View components (`LegacyTagView`, `LegacyCellView`, `LegacyButton`) plus their
`attrs.xml`, `colors.xml`, `dimens.xml` and one existing legacy screen layout
(`res/layout/activity_order_list_legacy.xml`).

**Not compiled here** (no Android SDK in this environment; `dl.google.com` / `maven.google.com` are blocked),
and intentionally not included in `settings.gradle.kts`. It exists to demonstrate View asset scanning:
a static scanner can read the Kotlin class declarations, `<declare-styleable>` enums/booleans/strings and the
custom-view usages in layouts — including one raw `<Button>` with a hard-coded `textColor` / `textSize`
that an audit should flag.

Target of the real app: compileSdk 32 / targetSdk 32 / minSdk 23 / buildTools 32.0.0.

| View | Figma component | Attributes |
|---|---|---|
| `com.demo.dsview.LegacyTagView` (AppCompatTextView) | Tag (Light only) | `tagTone` = primary / success / warning / danger |
| `com.demo.dsview.LegacyCellView` (LinearLayout) | Cell | `cellTitle`, `cellValue`, `cellShowArrow` |
| `com.demo.dsview.LegacyButton` (AppCompatButton) | Button | `btnType` = primary / default / text, `btnSize` = large / medium / small |
