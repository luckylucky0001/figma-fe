import { useEffect, useState } from 'react';

/**
 * CountDown — "剩余 HH:MM:SS", font.caption, color.danger (design/SPEC.md §2 CountDown).
 * Not yet in @ds (the spec marks it as Figma-only), so it lives with the page for now.
 */
export function CountDown({ seconds }: { seconds: number }) {
  const [left, setLeft] = useState(Math.max(0, Math.floor(seconds)));

  useEffect(() => {
    if (left <= 0) return;
    const timer = window.setTimeout(() => setLeft((s) => Math.max(0, s - 1)), 1000);
    return () => window.clearTimeout(timer);
  }, [left]);

  const pad = (n: number) => String(n).padStart(2, '0');
  const text = `${pad(Math.floor(left / 3600))}:${pad(Math.floor((left % 3600) / 60))}:${pad(left % 60)}`;

  return (
    <span
      data-figma="CountDown"
      style={{
        flexShrink: 0,
        whiteSpace: 'nowrap',
        fontSize: 'var(--ds-font-caption-size)',
        lineHeight: 'var(--ds-font-caption-line-height)',
        fontWeight: 'var(--ds-font-caption-weight)',
        color: 'var(--ds-color-danger)',
      }}
    >
      剩余 {text}
    </span>
  );
}
