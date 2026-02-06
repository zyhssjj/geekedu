
// 完整的 API 响应通用结构
export type ApiResponse<T> = {
    code: number;
    msg: string;
    data: T;
  };
  
  // 1. monthlySales API 的数据结构
  // 月度销售趋势数据中的每个系列
  export type MonthlySalesSeries = {
    name: string; // 例如: "2025年"
    data: number[]; // 对应月份的销售额数组
  };
  
  // monthlySales API 响应中 "data" 字段的结构
  export type MonthlySalesData = {
    months: string[]; // 月份名称数组
    series: MonthlySalesSeries[]; // 包含所有年度销售数据的系列数组
  };
  
  
  // 2. productStats API 的数据结构 (年龄分布)
  export type AgeDistributionItem = {
    name: string; 
    value: number; // 对应的百分比或计数
  };
  
  export type ProductStatsData = {
    pieData: AgeDistributionItem[];
  };
  
  // 3. trend API 的数据结构 (区域分布)
  export type RegionDistributionItem = {
    name: string; 
    value: number; // 对应的销售额
  };
  
  export type TrendData = { 
    regionData: RegionDistributionItem[];
  };
  
  
  // 4. tableList API 的数据结构 (渠道销售)
  
  export type ChannelSalesItem = {
    channel: string; // 例如: "WPS官网"
    sales: number; // 销售额
    orders: number; // 订单量
    avgPrice: number; // 平均价格
  };
  // 'asc' 升序, 'desc' 降序, null 表示未排序或默认
  export type SortOrder = 'asc' | 'desc' | null; // 'asc' 升序, 'desc' 降序, null 表示未排序或默认
