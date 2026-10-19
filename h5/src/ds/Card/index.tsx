import type { HTMLAttributes, ReactNode } from 'react';
import { cx } from '../utils';
import './index.css';

export interface CardProps extends Omit<HTMLAttributes<HTMLDivElement>, 'title'> {
  /** Optional card title (font.subtitle), 12px above the content. */
  title?: string;
  children?: ReactNode;
}

/**
 * Card — white container with 12px radius and 16px padding. Children stack vertically
 * with no gap; an optional title sits 12px above them.
 *
 * @example
 * <Card title="收货信息">
 *   <Cell title="收货人" value="张三" />
 * </Card>
 */
export function Card({ title, children, className, ...rest }: CardProps) {
  return (
    <div data-ds="Card" className={cx('ds-card', className)} {...rest}>
      {title !== undefined && title !== '' && <div className="ds-card__title">{title}</div>}
      <div className="ds-card__body">{children}</div>
    </div>
  );
}
