import type { HTMLAttributes } from 'react';
import { cx } from '../utils';
import './index.css';

export interface NavBarProps extends Omit<HTMLAttributes<HTMLDivElement>, 'title'> {
  /** Centered title text (font.title). */
  title: string;
  /** Called when the back chevron is tapped. */
  onBack?: () => void;
}

/**
 * NavBar — top navigation bar, 44px high, white background, back chevron on the left
 * (24×24, 12px from the left edge) and a centered title.
 *
 * @example
 * <NavBar title="订单详情" onBack={() => history.back()} />
 */
export function NavBar({ title, onBack, className, ...rest }: NavBarProps) {
  return (
    <div data-ds="NavBar" className={cx('ds-nav-bar', className)} {...rest}>
      <button type="button" className="ds-nav-bar__back" aria-label="返回" onClick={onBack}>
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" aria-hidden="true">
          <path d="M15 5L8 12L15 19" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
        </svg>
      </button>
      <div className="ds-nav-bar__title">{title}</div>
    </div>
  );
}
