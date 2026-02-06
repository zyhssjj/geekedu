// src/components/SalesTable.tsx
import type React from 'react';
import Th from './Th';
import type { ChannelSalesItem, SortOrder } from '../types'; // 导入 ChannelSalesItem 和 SortOrder
// 定义组件的 props 类型
type SalesTableProps = {
  data: ChannelSalesItem[]; // 期望接收 ChannelSalesItem[] 类型的详细数据
  sortKey: string | null; // 当前排序的列
  sortOrder: SortOrder; // 当前排序方向
  onSort: (key: string) => void; // 新增：排序处理函数
};
// 销售表格组件
const SalesTable: React.FC<SalesTableProps> = ({ data, sortKey, sortOrder, onSort }) => {
  return (
    <div className="overflow-x-auto rounded-lg shadow-lg">
      <table className="min-w-full divide-y divide-blue-700/50">
        <thead className="bg-blue-900/50">
          <tr>
            {/* 销售渠道：不可排序 */}
            <Th label="销售渠道" k="channel" sortable={false} currentSortKey={sortKey} currentSortOrder={sortOrder} onSort={onSort} />
            {/* 销售额：可排序 */}
            <Th label="销售额" k="sales" currentSortKey={sortKey} currentSortOrder={sortOrder} onSort={onSort} />
            {/* 订单量：可排序 */}
            <Th label="订单量" k="orders" currentSortKey={sortKey} currentSortOrder={sortOrder} onSort={onSort} />
            {/* 平均客单价：可排序 */}
            <Th label="平均客单价" k="avgPrice" currentSortKey={sortKey} currentSortOrder={sortOrder} onSort={onSort} />
          </tr>
        </thead>
        <tbody className="divide-y divide-blue-900/30">
          {data.length === 0 ? (
            <tr>
              <td colSpan={4} className="px-6 py-4 text-center text-sm text-blue-400">
                暂无渠道销售数据
              </td>
            </tr>
          ) : (
            data.map((row) => (
              <tr key={row.channel} className="hover:bg-blue-800/20 transition-colors duration-200">
                <td className="px-6 py-4 whitespace-nowrap text-sm text-blue-200">{row.channel}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-blue-200">¥{row.sales.toFixed(2)}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-blue-200">{row.orders}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-blue-200">{row.avgPrice.toFixed(2)}</td>
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  );
};

export default SalesTable;