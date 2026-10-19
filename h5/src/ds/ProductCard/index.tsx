import type { HTMLAttributes } from 'react';
import { Price } from '../Price';
import { cx } from '../utils';
import './index.css';

export interface ProductCardProps extends Omit<HTMLAttributes<HTMLDivElement>, 'title'> {
  /** Product name, clamped to two lines. */
  title: string;
  /** Spec / SKU description (font.caption, tertiary). */
  spec?: string;
  /** Unit price, rendered as `Price size="medium" tone="default"`. */
  price: number;
  /** Quantity, rendered as "x{count}". */
  count?: number;
  /** Thumbnail URL; a placeholder background shows when omitted. */
  thumb?: string;
}

/**
 * ProductCard — 80×80 thumbnail + info column (title, spec, bottom row with price and quantity).
 *
 * @example
 * <ProductCard title="无线降噪耳机 Pro" spec="曜石黑 / 标准版" price={899} count={1} />
 */
export function ProductCard({ title, spec, price, count, thumb, className, ...rest }: ProductCardProps) {
  return (
    <div data-ds="ProductCard" className={cx('ds-product-card', className)} {...rest}>
      <div className="ds-product-card__thumb">
        {thumb && <img src={thumb} alt="" />}
      </div>
      <div className="ds-product-card__info">
        <div className="ds-product-card__title">{title}</div>
        {spec !== undefined && spec !== '' && <div className="ds-product-card__spec">{spec}</div>}
        <div className="ds-product-card__bottom">
          <Price value={price} size="medium" tone="default" />
          {count !== undefined && <span className="ds-product-card__count">x{count}</span>}
        </div>
      </div>
    </div>
  );
}
