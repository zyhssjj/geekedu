// src/components/Th.tsx
import type React from 'react';
import type { SortOrder } from '../types'; 

type ThProps = {
  label: string;
  k: string; // 用于标识这一列的key，例如 'sales', 'orders', 'avgPrice'
  sortable?: boolean; // 是否可排序，默认为 true
  currentSortKey: string | null; // 当前排序的列
  currentSortOrder: SortOrder; // 当前排序方向
  onSort: (key: string) => void; // 排序处理函数
};

const Th: React.FC<ThProps> = ({ label, k, sortable = true, currentSortKey, currentSortOrder, onSort }) => {
  const isCurrentSortKey = currentSortKey === k;

  return (
    <th
      className={`px-6 py-3 text-left text-xs font-medium text-blue-300 uppercase tracking-wider ${
        sortable ? 'cursor-pointer hover:bg-blue-800/10 transition-colors duration-150' : ''
      }`}
      onClick={sortable ? () => onSort(k) : undefined}
    >
      <div className="flex items-center space-x-1">
        <span>{label}</span>
        {sortable && ( // 仅当列可排序时才显示图标
          <span className="ml-1 text-base">
            {isCurrentSortKey ? (
              currentSortOrder === 'asc' ? '▲' : '▼' // 当前排序，显示对应方向
            ) : (
              '⇅' // 可排序但不是当前排序，显示上下箭头
            )}
          </span>
        )}
      </div>
    </th>
  );
};

export default Th;