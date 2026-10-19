import type { HTMLAttributes, ReactNode } from 'react';
import { cx } from '../utils';
import './index.css';

export interface CellProps extends Omit<HTMLAttributes<HTMLDivElement>, 'title'> {
  /** Left title (font.body). */
  title: string;
  /** Optional sub-title under the title (font.caption, tertiary). */
  label?: string;
  /** Right-aligned value (font.body, secondary). */
  value?: ReactNode;
  /** Show a 16×16 right chevron after the value. */
  arrow?: boolean;
}

/**
 * Cell — list row: title (+ optional label) on the left, value (+ optional arrow) on the right.
 * 12px vertical padding, no horizontal padding (the Card provides it), min height 44.
 *
 * @example
 * <Cell title="配送方式" label="预计明天送达" value="快递" arrow />
 */
export function Cell({ title, label, value, arrow = false, className, ...rest }: CellProps) {
  return (
    <div data-ds="Cell" className={cx('ds-cell', className)} {...rest}>
      <div className="ds-cell__left">
        <div className="ds-cell__title">{title}</div>
        {label !== undefined && label !== '' && <div className="ds-cell__label">{label}</div>}
      </div>
      {(value !== undefined || arrow) && (
        <div className="ds-cell__right">
          {value !== undefined && <div className="ds-cell__value">{value}</div>}
          {arrow && (
            <svg className="ds-cell__arrow" width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden="true">
              <path d="M6 3.5L10.5 8L6 12.5" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
            </svg>
          )}
        </div>
      )}
    </div>
  );
}
