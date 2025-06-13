// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

import tailwindcss from '@tailwindcss/vite';

// https://astro.build/config
export default defineConfig({
  integrations: [
      starlight({
          title: 'Koito',
          social: [{ icon: 'github', label: 'GitHub', href: 'https://github.com/gabehf/koito' }],
          sidebar: [
              {
                  label: 'Guides',
                  items: [
                      // Each item here is one entry in the navigation menu.
                      { label: 'Installation', slug: 'guides/installation' },
                      { label: 'Importing Data', slug: 'guides/importing' },
                      { label: 'Setting up the Scrobbler', slug: 'guides/scrobbler' },
                      { label: 'Editing Data', slug: 'guides/editing' },
                  ],
              },
              {
                  label: 'Reference',
                  items: [
                    { label: 'Configuration Options', slug: 'reference/configuration' },
                  ]
              },
          ],
		  customCss: [
			// Path to your Tailwind base styles:
			'./src/styles/global.css',
		  ],
      }),
	],

  vite: {
    plugins: [tailwindcss()],
  },
});