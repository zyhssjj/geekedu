// src/App.tsx
import type React from 'react';
import { useEffect, useState, useMemo } from 'react'; // 引入 useMemo 优化排序
import axios from 'axios';
import { SalesTrend } from './components/SalesTrend';
import { AgeDist, RegionDist } from './components/OtherCharts';
import SalesTable from './components/SalesTable';
// 导入所有需要的类型
import type {
  ApiResponse,
  MonthlySalesData,
  MonthlySalesSeries,
  AgeDistributionItem,
  ProductStatsData,
  RegionDistributionItem,
  TrendData,
  ChannelSalesItem,
  SortOrder // 导入 SortOrder 类型
} from './types';

// API 地址配置
const API_BASE_URL = 'https://m1.apifoxmock.com/m1/5076419-0-default';
const API_TOKEN = 'kS5RF-4neMuhg_Qvabu40'; 

const API_URLS = {
  monthlySales: `${API_BASE_URL}/api/sales/monthlySales?apifoxToken=${API_TOKEN}`,
  productStats: `${API_BASE_URL}/api/sales/productStats?apifoxToken=${API_TOKEN}`,
  trend: `${API_BASE_URL}/api/sales/trend?apifoxToken=${API_TOKEN}`,
  tableList: `${API_BASE_URL}/api/sales/tableList?apifoxToken=${API_TOKEN}`,
};
// 本地存储排序状态的键名
const LOCAL_STORAGE_SORT_KEY = 'channelSalesSortState'; 

const App: React.FC = () => {
  // ----------------- State for Monthly Sales Trend -----------------
  const [monthlySalesMonths, setMonthlySalesMonths] = useState<string[]>([]);
  const [monthlySalesSeries, setMonthlySalesSeries] = useState<MonthlySalesSeries[]>([]);
  const [monthlySalesLoading, setMonthlySalesLoading] = useState(true);
  const [monthlySalesError, setMonthlySalesError] = useState<string | null>(null);

  // ----------------- State for Product Stats (Age Distribution) -----------------
  const [ageDistributionData, setAgeDistributionData] = useState<AgeDistributionItem[]>([]);
  const [productStatsLoading, setProductStatsLoading] = useState(true);
  const [productStatsError, setProductStatsError] = useState<string | null>(null);

  // ----------------- State for Trend (Region Distribution) -----------------
  const [regionDistributionData, setRegionDistributionData] = useState<RegionDistributionItem[]>([]);
  const [regionTrendLoading, setRegionTrendLoading] = useState(true);
  const [regionTrendError, setRegionTrendError] = useState<string | null>(null);

  // ----------------- State for Channel Sales Table (tableList) -----------------
  const [channelSalesData, setChannelSalesData] = useState<ChannelSalesItem[]>([]);
  const [channelSalesLoading, setChannelSalesLoading] = useState(true);
  const [channelSalesError, setChannelSalesError] = useState<string | null>(null);

  // 表格排序状态
  const [sortKey, setSortKey] = useState<string | null>(null);
  const [sortOrder, setSortOrder] = useState<SortOrder>(null);


  // Helper function for API calls 
  const fetchData = async <T,>(
    url: string,
    onSuccess: (data: T) => void,
    setLoading: React.Dispatch<React.SetStateAction<boolean>>,
    setError: React.Dispatch<React.SetStateAction<string | null>>,
  ) => {
    setLoading(true);
    setError(null);
    try {
      console.log(`🔵 正在从 API 获取数据: ${url}`);
      const response = await axios.get<ApiResponse<any>>(url);

      console.log(`🟢 【API 原始响应】URL: ${url}`, response.data);

      if (response.data.code === 0 && response.data.data !== undefined && response.data.data !== null) {

    
        onSuccess(response.data.data as T);
        console.log(`✅ 从 ${url} 提取到数据。`);
      } else {
        console.error(`❌ API 返回失败或数据结构异常。URL: ${url}, 响应:`, response.data);
        throw new Error(`API 返回失败或数据结构异常: ${response.data.msg || '未知错误'}`);
      }
    } catch (err) {
      console.error(`🔴 数据获取失败捕获异常。URL: ${url}`, err);
      if (axios.isAxiosError(err)) {
        if (err.response) {
            setError(`错误: ${err.response.status} - ${err.response.statusText || '服务器响应错误。'}`);
            console.error("响应数据:", err.response.data);
        } else if (err.request) {
            setError("错误: 未收到服务器响应。请检查网络或 API 服务器状态。");
            console.error("请求数据:", err.request);
        } else {
            setError(`错误: ${err.message}`);
        }
      } else if (err instanceof Error) {
        setError(`错误: ${err.message}`);
      } else {
        setError(`未知错误: ${String(err)}`);
      }
    } finally {
      setLoading(false);
    }
  };

  // Effect for Monthly Sales Trend
  useEffect(() => {
    fetchData<MonthlySalesData>(
      API_URLS.monthlySales,
      (data) => {
        setMonthlySalesMonths(data.months);
        setMonthlySalesSeries(data.series);
      },
      setMonthlySalesLoading,
      setMonthlySalesError
    );
  }, []);

  // Effect for Product Stats (Age Distribution)
  useEffect(() => {
    fetchData<ProductStatsData>(
      API_URLS.productStats,
      (data) => setAgeDistributionData(data.pieData), // 从 ProductStatsData 中提取 pieData
      setProductStatsLoading,
      setProductStatsError
    );
  }, []);

  // Effect for Trend (Region Distribution)
  useEffect(() => {
    fetchData<TrendData>(
      API_URLS.trend,
      (data) => setRegionDistributionData(data.regionData), // 从 TrendData 中提取 regionData
      setRegionTrendLoading,
      setRegionTrendError
    );
  }, []);

  // Effect for Channel Sales Table (tableList)
  useEffect(() => {
    fetchData<ChannelSalesItem[]>(
      API_URLS.tableList,
      setChannelSalesData, 
      setChannelSalesLoading,
      setChannelSalesError
    );
  }, []);

  // 从 localStorage 加载排序状态
  useEffect(() => {
    try {
      const savedSortState = localStorage.getItem(LOCAL_STORAGE_SORT_KEY);
      if (savedSortState) {
        const { key, order } = JSON.parse(savedSortState);
        setSortKey(key);
        setSortOrder(order);
      }
    } catch (e) {
      console.error("从 localStorage 加载排序状态失败:", e);
    // 如果解析失败，清除旧的无效数据
      localStorage.removeItem(LOCAL_STORAGE_SORT_KEY);
    }
  }, []); // 仅在组件挂载时执行一次

  // 当排序状态改变时，保存到 localStorage
  useEffect(() => {
    if (sortKey !== null && sortOrder !== null) {
      localStorage.setItem(LOCAL_STORAGE_SORT_KEY, JSON.stringify({ key: sortKey, order: sortOrder }));
    } else {
      // 如果没有排序，或排序被重置，则清除 localStorage 中的排序记录
      localStorage.removeItem(LOCAL_STORAGE_SORT_KEY);
    }
  }, [sortKey, sortOrder]); // 依赖于 sortKey 和 sortOrder 的变化

  // 处理表格排序的函数
  const handleSort = (key: string) => {
    if (sortKey === key) {
      // 如果点击的是当前排序的列，则切换排序方向
      setSortOrder(prevOrder => {
        if (prevOrder === 'asc') return 'desc';
        if (prevOrder === 'desc') return null; // 第三次点击取消排序 (或可以改为 'asc' 循环)
        return 'asc'; // 默认升序
      });
    } else {
      // 如果点击的是新的列，则将该列设为升序排序
      setSortKey(key);
      setSortOrder('asc');
    }
  };

  // 使用 useMemo 对 channelSalesData 进行排序，只有在数据或排序状态改变时才重新计算
  const sortedChannelSalesData = useMemo(() => {
    if (!sortKey || !sortOrder || channelSalesData.length === 0) {
      return channelSalesData; // 没有排序键或数据为空，返回原始数据
    }

    const sortedData = [...channelSalesData].sort((a, b) => {
      // 这里的 key 需要与 ChannelSalesItem 的属性匹配
      const aValue = a[sortKey as keyof ChannelSalesItem];
      const bValue = b[sortKey as keyof ChannelSalesItem];

      // 确保是数字类型进行比较
      if (typeof aValue === 'number' && typeof bValue === 'number') {
        return sortOrder === 'asc' ? aValue - bValue : bValue - aValue;
      }
      // 如果未来有其他非数字的可排序列，可以在这里添加更多比较逻辑
      return 0; // 默认不改变顺序
    });
    return sortedData;
  }, [channelSalesData, sortKey, sortOrder]);


  // Combine loading and error states for main UI rendering
  const overallLoading = monthlySalesLoading || productStatsLoading || regionTrendLoading || channelSalesLoading;
  const overallError = monthlySalesError || productStatsError || regionTrendError || channelSalesError;


  return (
    <div className="min-h-screen bg-[#060b28] p-6 text-white">
      {/* 标题 */}
      <header className="mb-8 relative text-center">
        <h1 className="text-3xl font-bold tracking-[0.2em] text-transparent bg-clip-text bg-gradient-to-b from-cyan-300 to-blue-600">
          WPS会员销售数据仪表盘
        </h1>
        <div className="mt-2 h-[2px] w-[400px] mx-auto bg-gradient-to-r from-transparent via-blue-500 to-transparent"></div>
      </header>

      {/* 加载/错误提示 */}
      {overallLoading && (
        <div className="text-center text-lg mt-8">正在加载所有数据...</div>
      )}
      {overallError && (
        <div className="text-center text-red-500 text-lg mt-8">
          数据加载错误: <br/>
          {monthlySalesError && <div>月度销售: {monthlySalesError}</div>}
          {productStatsError && <div>年龄分布: {productStatsError}</div>}
          {regionTrendError && <div>区域分布: {regionTrendError}</div>}
          {channelSalesError && <div>渠道销售: {channelSalesError}</div>}
        </div>
      )}

      {/* 栅格布局 - 仅在所有数据加载完成且无错误时显示 */}
      {!overallLoading && !overallError && (
        <main className="max-w-[1600px] mx-auto grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* 月度销售趋势图 */}
          <SalesTrend months={monthlySalesMonths} series={monthlySalesSeries} />
          {/* 用户年龄分布图 */}
          <AgeDist ageData={ageDistributionData} />
          {/* 区域销售额分布图 */}
          <RegionDist regionData={regionDistributionData} />
          {/* 渠道销售表格 - 传递排序相关 props */}
          <SalesTable
            data={sortedChannelSalesData}
            sortKey={sortKey}
            sortOrder={sortOrder}
            onSort={handleSort}
          />
        </main>
      )}
    </div>
  );
};

export default App;