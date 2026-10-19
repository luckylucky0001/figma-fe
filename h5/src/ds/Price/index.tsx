import type { HTMLAttributes } from 'react';
import { cx } from '../utils';
import './index.css';

export interface PriceProps extends Omit<HTMLAttributes<HTMLSpanElement>, 'children'> {
  /** Amount in yuan, e.g. `128.5` renders as ¥128.50. */
  value: number;
  /** Large: ¥14 / 22 / 14 · Medium: ¥12 / 18 / 12 · Small: ¥12 / 14 / 12. Default `'medium'`. */
  size?: 'large' | 'medium' | 'small';
  /** `'primary'` uses color.primary, `'default'` uses text-primary. Default `'primary'`. */
  tone?: 'primary' | 'default';
}

/**
 * Price — "¥" + integer + "." + two decimals, baseline aligned. Integer is 700 weight,
 * symbol and decimals are 600.
 *
 * @example
 * <Price value={299} size="large" />
 * <Price value={59.9} size="medium" tone="default" />
 */
export function Price({ value, size = 'medium', tone = 'primary', className, ...rest }: PriceProps) {
  const fixed = (Number.isFinite(value) ? Math.abs(value) : 0).toFixed(2);
  const [integer, decimal] = fixed.split('.');
  const negative = value < 0;
  return (
    <span
      data-ds="Price"
      className={cx('ds-price', `ds-price--${size}`, `ds-price--${tone}`, className)}
      {...rest}
    >
      <span className="ds-price__symbol">{negative ? '-¥' : '¥'}</span>
      <span className="ds-price__integer">{integer}</span>
      <span className="ds-price__decimal">.{decimal}</span>
    </span>
  );
}
