import type { ButtonHTMLAttributes, MouseEventHandler, ReactNode } from 'react';
import { cx } from '../utils';
import './index.css';

export interface ButtonProps extends Omit<ButtonHTMLAttributes<HTMLButtonElement>, 'type'> {
  /** Visual type. Default `'default'`. */
  type?: 'primary' | 'default' | 'text';
  /** Large h44 · Medium h32 · Small h24. Default `'medium'`. */
  size?: 'large' | 'medium' | 'small';
  /** Stretch to the full width of the container. */
  block?: boolean;
  onClick?: MouseEventHandler<HTMLButtonElement>;
  children?: ReactNode;
}

/**
 * Button — pill button. Primary (brand fill, white text), Default (white fill, 1px border),
 * Text (no fill/border, brand text). Sizes Large 44 / Medium 32 / Small 24 px high.
 *
 * @example
 * <Button type="primary" size="medium" onClick={pay}>去支付</Button>
 * <Button>取消订单</Button>
 */
export function Button({
  type = 'default',
  size = 'medium',
  block = false,
  onClick,
  children,
  className,
  ...rest
}: ButtonProps) {
  return (
    <button
      type="button"
      data-ds="Button"
      className={cx('ds-button', `ds-button--${type}`, `ds-button--${size}`, block && 'ds-button--block', className)}
      onClick={onClick}
      {...rest}
    >
      {children}
    </button>
  );
}
