// src/components/SalesChart.tsx
import type React from 'react';
import ReactECharts from 'echarts-for-react';
import type { EChartsOption } from 'echarts'; // 导入 EChartsOption 类型
// 定义组件的 props 类型
type SalesChartProps = {
  option: EChartsOption; // ECharts配置对象
  className?: string; // 用于Tailwind样式
};
// 销售图表组件
const SalesChart: React.FC<SalesChartProps> = ({ option, className }) => {
  return (
    <ReactECharts
      option={option}
      className={className}
      notMerge={true}
      lazyUpdate={true}
    />
  );
};

export default SalesChart;