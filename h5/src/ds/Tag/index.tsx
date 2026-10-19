import type { HTMLAttributes, ReactNode } from 'react';
import { cx } from '../utils';
import './index.css';

export interface TagProps extends HTMLAttributes<HTMLSpanElement> {
  /** Tone. Default `'primary'`. */
  type?: 'primary' | 'success' | 'warning' | 'danger';
  /** `true` = Light variant (tinted bg, toned text); `false` = Solid (toned bg, white text). Default `true`. */
  light?: boolean;
  children?: ReactNode;
}

/**
 * Tag — small status label, 20px high, 6px horizontal padding, 4px radius, 12/20/500 text.
 *
 * @example
 * <Tag type="success">已发货</Tag>
 * <Tag type="danger" light={false}>超时</Tag>
 */
export function Tag({ type = 'primary', light = true, children, className, ...rest }: TagProps) {
  return (
    <span
      data-ds="Tag"
      className={cx('ds-tag', `ds-tag--${type}`, light ? 'ds-tag--light' : 'ds-tag--solid', className)}
      {...rest}
    >
      {children}
    </span>
  );
}
