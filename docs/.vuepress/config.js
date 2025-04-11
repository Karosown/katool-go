module.exports = {
  title: 'Katool-Go 文档',
  description: 'Katool-Go 是一个功能丰富的 Go 工具库，借鉴了 Java 生态中的优秀设计',
  base: '/katool-go/',
  themeConfig: {
    nav: [
      { text: '首页', link: '/' },
      { text: 'GitHub', link: 'https://github.com/karosown/katool-go' },
    ],
    sidebar: [
      {
        title: '入门指南',
        collapsable: false,
        children: [
          '/',
          '/guide/getting-started',
        ]
      },
      {
        title: '核心功能',
        collapsable: false,
        children: [
          '/stream',
          '/lists',
          '/ioc',
          '/lock',
          '/convert',
        ]
      },
      {
        title: '工具模块',
        collapsable: false,
        children: [
          '/web_crawler',
          '/algorithm',
          '/log',
          '/file',
          '/db',
        ]
      },
    ],
    sidebarDepth: 2,
    displayAllHeaders: false,
    activeHeaderLinks: true,
  },
  plugins: [
    '@vuepress/back-to-top',
    '@vuepress/medium-zoom',
    [
      '@vuepress/search', 
      {
        searchMaxSuggestions: 10
      }
    ]
  ],
  markdown: {
    lineNumbers: true
  }
} 