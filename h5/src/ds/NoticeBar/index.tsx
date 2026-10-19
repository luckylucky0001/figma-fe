import type { HTMLAttributes, ReactNode } from 'react';
import { cx } from '../utils';
import './index.css';

export interface NoticeBarProps extends HTMLAttributes<HTMLDivElement> {
  /** `'info'` = primary-light bg / primary text; `'warning'` = warning-light bg / warning text. Default `'info'`. */
  tone?: 'info' | 'warning';
  children?: ReactNode;
}

/**
 * NoticeBar — single-line notice, 36px high, 12px horizontal padding, 8px radius, font.caption.
 * Overflowing text is truncated with an ellipsis.
 *
 * @example
 * <NoticeBar tone="warning">受天气影响，部分地区配送可能延迟</NoticeBar>
 */
export function NoticeBar({ tone = 'info', children, className, ...rest }: NoticeBarProps) {
  return (
    <div data-ds="NoticeBar" className={cx('ds-notice-bar', `ds-notice-bar--${tone}`, className)} {...rest}>
      <span className="ds-notice-bar__text">{children}</span>
    </div>
  );
}
