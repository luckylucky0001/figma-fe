import { StrictMode, type ComponentType } from 'react';
import { createRoot } from 'react-dom/client';
import '@ds/tokens.css';
import '@ds/global.css';

type PageModule = { default: ComponentType };

const modules = import.meta.glob<PageModule>('./pages/*/index.tsx', { eager: true });

const pages: Record<string, ComponentType> = {};
for (const [path, mod] of Object.entries(modules)) {
  const name = path.match(/^\.\/pages\/([^/]+)\/index\.tsx$/)?.[1];
  if (name && mod.default) pages[name] = mod.default;
}

function PageList({ missing }: { missing: string | null }) {
  return (
    <div style={{ padding: 16, fontSize: 14, lineHeight: '22px' }}>
      {missing && <p style={{ margin: '0 0 12px' }}>Page “{missing}” not found.</p>}
      <p style={{ margin: '0 0 8px', fontWeight: 600 }}>Pages</p>
      <ul style={{ margin: 0, paddingLeft: 20 }}>
        {Object.keys(pages)
          .sort()
          .map((name) => (
            <li key={name}>
              <a href={`?page=${encodeURIComponent(name)}`}>{name}</a>
            </li>
          ))}
      </ul>
    </div>
  );
}

const pageName = new URLSearchParams(window.location.search).get('page');
const Page = pageName ? pages[pageName] : undefined;

createRoot(document.getElementById('root')!).render(
  <StrictMode>{Page ? <Page /> : <PageList missing={pageName} />}</StrictMode>,
);
