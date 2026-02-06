// src/components/SalesTrend.tsx
import type React from 'react';
import SalesChart from './SalesChart';
import type { EChartsOption, SeriesOption } from 'echarts';
import type { MonthlySalesSeries } from '../types';

type SalesTrendProps = {
  months: string[]; // 例如: ["1月", "2月", ...]
  series: MonthlySalesSeries[]; // 例如: [{ name: "2025年", data: [...] }, ...]
};
// 月度销售趋势组件
export const SalesTrend: React.FC<SalesTrendProps> = ({ months, series }) => {
  const echartsSeries: SeriesOption[] = series.map(s => ({
    name: s.name,
    type: 'line', 
    data: s.data,
    focus: 'series',
    smooth: true, 
  }));

  const chartOption: EChartsOption = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    legend: {
      data: series.map(s => s.name),
      textStyle: {
        color: '#B3E0FF'
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: months,
      axisLabel: {
        color: '#B3E0FF'
      }
    },
    yAxis: {
      type: 'value',
      name: '销售额 (¥)',
      nameTextStyle: {
        color: '#B3E0FF'
      },
      axisLabel: {
        color: '#B3E0FF'
      },
      splitLine: {
        lineStyle: {
          color: 'rgba(255,255,255,0.1)'
        }
      }
    },
    series: echartsSeries
  };

  return (
    <div className="bg-blue-900 p-6 rounded-lg shadow-lg">
      <h2 className="text-xl font-semibold mb-4 text-blue-200">月度销售趋势 (多年对比)</h2>
      {months.length > 0 && series.length > 0 ? (
        <SalesChart option={chartOption} className="w-full h-80" />
      ) : (
        <div className="h-80 flex items-center justify-center text-blue-400">暂无月度销售趋势数据</div>
      )}
    </div>
  );
};