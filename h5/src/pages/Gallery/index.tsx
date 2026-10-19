import type { CSSProperties, ReactNode } from 'react';
import {
  ActionBar,
  Button,
  Card,
  Cell,
  Divider,
  NavBar,
  NoticeBar,
  Price,
  ProductCard,
  Tag,
} from '@ds';
import type { ButtonProps, PriceProps, TagProps } from '@ds';

const labelStyle: CSSProperties = {
  margin: '0 0 6px',
  fontSize: 'var(--ds-font-caption-size)',
  lineHeight: 'var(--ds-font-caption-line-height)',
  color: 'var(--ds-color-text-tertiary)',
};

const rowStyle: CSSProperties = { display: 'flex', flexWrap: 'wrap', alignItems: 'center', gap: 8 };

function Section({ name, children }: { name: string; children: ReactNode }) {
  return (
    <section data-gallery={name} style={{ padding: '16px 12px 0' }}>
      <h2
        style={{
          margin: '0 0 8px',
          fontSize: 'var(--ds-font-subtitle-size)',
          lineHeight: 'var(--ds-font-subtitle-line-height)',
          fontWeight: 600,
          color: 'var(--ds-color-text-primary)',
        }}
      >
        {name}
      </h2>
      <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>{children}</div>
    </section>
  );
}

function Variant({ label, children }: { label: string; children: ReactNode }) {
  return (
    <div>
      <p style={labelStyle}>{label}</p>
      {children}
    </div>
  );
}

const tagTypes: NonNullable<TagProps['type']>[] = ['primary', 'success', 'warning', 'danger'];
const buttonTypes: NonNullable<ButtonProps['type']>[] = ['primary', 'default', 'text'];
const sizes: NonNullable<ButtonProps['size']>[] = ['large', 'medium', 'small'];
const priceTones: NonNullable<PriceProps['tone']>[] = ['primary', 'default'];

/** Gallery — every component in every variant, for visual QA. */
export default function Gallery() {
  return (
    <div style={{ minHeight: '100vh', background: 'var(--ds-color-bg-page)', paddingBottom: 56 + 16 }}>
      <Section name="NavBar">
        <Variant label="default">
          <NavBar title="订单详情" onBack={() => {}} data-qa="navbar" />
        </Variant>
      </Section>

      <Section name="Card">
        <Variant label="Title=On">
          <Card title="收货信息" data-qa="card-title">
            <div style={{ fontSize: 14, lineHeight: '20px' }}>卡片内容区域</div>
          </Card>
        </Variant>
        <Variant label="Title=Off">
          <Card data-qa="card">
            <div style={{ fontSize: 14, lineHeight: '20px' }}>卡片内容区域</div>
          </Card>
        </Variant>
      </Section>

      <Section name="Cell">
        <Card>
          <Cell title="订单编号" value="DD20261018093021" data-qa="cell" />
          <Divider />
          <Cell title="配送方式" value="快递" arrow />
          <Divider />
          <Cell title="发票" label="电子普通发票 · 个人" value="已开具" />
          <Divider />
          <Cell title="优惠券" label="满 300 减 30" value="-¥30.00" arrow />
          <Divider />
          <Cell title="仅标题" />
        </Card>
        <p style={labelStyle}>Arrow=Off/Label=Off · Arrow=On · Label=On · Arrow=On+Label=On · title only</p>
      </Section>

      <Section name="Tag">
        <Variant label="Variant=Light">
          <div style={rowStyle}>
            {tagTypes.map((t) => (
              <Tag key={t} type={t} data-qa={`tag-light-${t}`}>
                {t === 'primary' ? '待付款' : t === 'success' ? '已发货' : t === 'warning' ? '待发货' : '已取消'}
              </Tag>
            ))}
          </div>
        </Variant>
        <Variant label="Variant=Solid">
          <div style={rowStyle}>
            {tagTypes.map((t) => (
              <Tag key={t} type={t} light={false}>
                {t}
              </Tag>
            ))}
          </div>
        </Variant>
      </Section>

      <Section name="Price">
        {priceTones.map((tone) => (
          <Variant key={tone} label={`Tone=${tone} · Large / Medium / Small`}>
            <div style={{ ...rowStyle, alignItems: 'baseline', gap: 16 }}>
              <Price value={1299} size="large" tone={tone} />
              <Price value={59.9} size="medium" tone={tone} />
              <Price value={8.05} size="small" tone={tone} />
            </div>
          </Variant>
        ))}
      </Section>

      <Section name="Button">
        {buttonTypes.map((type) => (
          <Variant key={type} label={`Type=${type} · Large / Medium / Small`}>
            <div style={rowStyle}>
              {sizes.map((size) => (
                <Button key={size} type={type} size={size} data-qa={`button-${type}-${size}`}>
                  {type === 'text' ? '查看详情' : '提醒发货'}
                </Button>
              ))}
            </div>
          </Variant>
        ))}
        <Variant label="block">
          <Button type="primary" size="large" block>
            立即支付
          </Button>
        </Variant>
        <Variant label="disabled">
          <Button type="primary" disabled>
            不可点击
          </Button>
        </Variant>
      </Section>

      <Section name="Divider">
        <Card>
          <div style={{ fontSize: 14, lineHeight: '20px', paddingBottom: 12 }}>上方内容</div>
          <Divider data-qa="divider" />
          <div style={{ fontSize: 14, lineHeight: '20px', paddingTop: 12 }}>下方内容</div>
        </Card>
      </Section>

      <Section name="ProductCard">
        <Card>
          <ProductCard
            title="无线降噪蓝牙耳机 Pro 主动降噪 长续航 入耳式运动耳机 超长标题测试两行省略"
            spec="曜石黑 / 标准版"
            price={899}
            count={1}
            data-qa="product-card"
          />
        </Card>
        <Card>
          <ProductCard title="北欧风陶瓷马克杯" price={39.9} count={2} />
        </Card>
      </Section>

      <Section name="NoticeBar">
        <Variant label="Tone=Info">
          <NoticeBar tone="info" data-qa="notice-bar">
            订单已提交，请在 30 分钟内完成支付
          </NoticeBar>
        </Variant>
        <Variant label="Tone=Warning (overflow)">
          <NoticeBar tone="warning">受天气影响，部分地区配送可能延迟 1-2 天，请耐心等待，给您带来不便敬请谅解</NoticeBar>
        </Variant>
      </Section>

      <Section name="ActionBar">
        <p style={labelStyle}>fixed to the bottom of the viewport ↓</p>
      </Section>

      <ActionBar data-qa="action-bar">
        <Button>取消订单</Button>
        <Button type="primary">去支付</Button>
      </ActionBar>
    </div>
  );
}
