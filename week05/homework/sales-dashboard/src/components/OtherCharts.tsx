// src/components/OtherCharts.tsx
import type React from 'react';
import SalesChart from './SalesChart';
import type { EChartsOption } from 'echarts';
import type { AgeDistributionItem, RegionDistributionItem } from '../types';

// AgeDist 组件，现在接收真实数据
type AgeDistProps = {
  ageData: AgeDistributionItem[];
};

export const AgeDist: React.FC<AgeDistProps> = ({ ageData }) => {
  const chartOption: EChartsOption = {
    tooltip: {
      trigger: 'item',
      formatter: '{b} : {c}%' // 假设数据是百分比
    },
    legend: {
      orient: 'vertical',
      left: 'left',
      textStyle: {
        color: '#B3E0FF'
      }
    },
    series: [
      {
        name: '用户年龄分布',
        type: 'pie',
        radius: '50%',
        data: ageData, // 直接使用 ageData
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        },
        label: {
          color: '#B3E0FF',
          formatter: '{b} ({d}%)' // 显示名称和百分比
        }
      }
    ]
  };

  return (
    <div className="bg-blue-900 p-6 rounded-lg shadow-lg">
      <h2 className="text-xl font-semibold mb-4 text-blue-200">用户年龄分布</h2>
      {ageData.length > 0 ? (
        <SalesChart option={chartOption} className="w-full h-80" />
      ) : (
        <div className="h-80 flex items-center justify-center text-blue-400">暂无年龄分布数据</div>
      )}
    </div>
  );
};

// RegionDist 组件，现在接收真实数据
type RegionDistProps = {
  regionData: RegionDistributionItem[];
};

export const RegionDist: React.FC<RegionDistProps> = ({ regionData }) => {
  const chartOption: EChartsOption = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      },
      formatter: '{b}: ¥{c}'
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: regionData.map(item => item.name), // 从数据中提取名称
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
    series: [
      {
        name: '区域销售额',
        type: 'bar',
        data: regionData.map(item => item.value), // 从数据中提取值
        itemStyle: {
          color: '#00FA9A'
        },
        label: {
          show: true,
          position: 'top',
          color: '#B3E0FF',
          formatter: '¥{c}'
        }
      }
    ]
  };

  return (
    <div className="bg-blue-900 p-6 rounded-lg shadow-lg">
      <h2 className="text-xl font-semibold mb-4 text-blue-200">区域销售额分布</h2>
      {regionData.length > 0 ? (
        <SalesChart option={chartOption} className="w-full h-80" />
      ) : (
        <div className="h-80 flex items-center justify-center text-blue-400">暂无区域销售数据</div>
      )}
    </div>
  );
};