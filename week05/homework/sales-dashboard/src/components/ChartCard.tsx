import React from 'react';

interface Props {
  title: string;
  children: React.ReactNode;
}

export const ChartCard: React.FC<Props> = ({ title, children }) => (
  <div className="bg-[#10163a] border border-blue-900/40 rounded-lg p-4 shadow-2xl relative">
    {/* 四角装饰 */}
    <div className="absolute top-0 left-0 w-3 h-3 border-t-2 border-l-2 border-cyan-500 rounded-tl-sm"></div>
    <div className="absolute top-0 right-0 w-3 h-3 border-t-2 border-r-2 border-cyan-500 rounded-tr-sm"></div>
    <div className="absolute bottom-0 left-0 w-3 h-3 border-b-2 border-l-2 border-cyan-500 rounded-bl-sm"></div>
    <div className="absolute bottom-0 right-0 w-3 h-3 border-b-2 border-r-2 border-cyan-500 rounded-br-sm"></div>
    
    <h3 className="text-center text-gray-200 font-bold mb-4 tracking-wider">{title}</h3>
    <div className="h-[320px] w-full">{children}</div>
  </div>
);