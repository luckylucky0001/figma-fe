import type { HTMLAttributes, ReactNode } from 'react';
import { cx } from '../utils';
import './index.css';

export interface ActionBarProps extends HTMLAttributes<HTMLDivElement> {
  /** Usually Buttons; right-aligned with an 8px gap. */
  children?: ReactNode;
}

/**
 * ActionBar — bar fixed to the bottom of the viewport, 56px high, white background,
 * 1px top border, 16px horizontal padding; children are right-aligned with an 8px gap.
 * Pages should reserve 56px of bottom padding so content is not covered.
 *
 * @example
 * <ActionBar>
 *   <Button>取消订单</Button>
 *   <Button type="primary">去支付</Button>
 * </ActionBar>
 */
export function ActionBar({ children, className, ...rest }: ActionBarProps) {
  return (
    <div data-ds="ActionBar" className={cx('ds-action-bar', className)} {...rest}>
      {children}
    </div>
  );
}
