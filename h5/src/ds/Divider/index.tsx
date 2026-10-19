import type { HTMLAttributes } from 'react';
import { cx } from '../utils';
import './index.css';

/** Divider has no own props; only passthrough attributes (data-*, className, style). */
// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface DividerProps extends Omit<HTMLAttributes<HTMLDivElement>, 'children'> {}

/**
 * Divider — 1px horizontal line in color.border, no margin.
 *
 * @example
 * <Divider />
 */
export function Divider({ className, ...rest }: DividerProps) {
  return <div data-ds="Divider" role="separator" className={cx('ds-divider', className)} {...rest} />;
}
