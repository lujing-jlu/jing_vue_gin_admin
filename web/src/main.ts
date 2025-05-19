import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import 'element-plus/dist/index.css'
import './style.css'
import App from './App.vue'
import router from './router'

// 引入ECharts
import ECharts from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { BarChart, LineChart, PieChart, GaugeChart } from 'echarts/charts'
import { 
  TitleComponent, TooltipComponent, LegendComponent, GridComponent,
  DatasetComponent, ToolboxComponent
} from 'echarts/components'

// 全局注册ECharts组件
use([
  CanvasRenderer,
  BarChart, LineChart, PieChart, GaugeChart,
  TitleComponent, TooltipComponent, LegendComponent, GridComponent,
  DatasetComponent, ToolboxComponent
])

const app = createApp(App)

// 全局注册ElementPlus图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// 全局注册ECharts组件
app.component('v-chart', ECharts)

app.use(createPinia())
app.use(router)
app.use(ElementPlus)
app.mount('#app')
