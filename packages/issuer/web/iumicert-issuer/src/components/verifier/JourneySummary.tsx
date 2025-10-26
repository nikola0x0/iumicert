import { JourneyReceipt } from "./types";

interface JourneySummaryProps {
  receipt: JourneyReceipt;
}

export function JourneySummary({ receipt }: JourneySummaryProps) {
  return (
    <div className="bg-white rounded-2xl shadow-lg border border-gray-100 p-6">
      <h2 className="text-xl font-bold text-slate-900 mb-4">
        Journey Summary
      </h2>
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="bg-blue-50 p-5 rounded-xl border border-blue-100">
          <div className="text-2xl font-bold text-blue-600">
            {Object.keys(receipt.term_receipts).length}
          </div>
          <div className="text-sm text-gray-600 mt-1">Terms Completed</div>
        </div>
        <div className="bg-green-50 p-5 rounded-xl border border-green-100">
          <div className="text-2xl font-bold text-green-600">
            {Object.values(receipt.term_receipts).reduce(
              (sum, term) => sum + term.total_courses,
              0
            )}
          </div>
          <div className="text-sm text-gray-600 mt-1">Total Courses</div>
        </div>
        <div className="bg-purple-50 p-5 rounded-xl border border-purple-100">
          <div className="text-2xl font-bold text-purple-600">
            {receipt.blockchain_ready ? "Ready" : "Pending"}
          </div>
          <div className="text-sm text-gray-600 mt-1">Blockchain Status</div>
        </div>
        <div className="bg-indigo-50 p-5 rounded-xl border border-indigo-100">
          <div className="text-sm font-bold text-indigo-600">
            {receipt.receipt_type.selective_disclosure ? "Selective" : "Full"}
          </div>
          <div className="text-sm text-gray-600 mt-1">Disclosure Type</div>
        </div>
      </div>
    </div>
  );
}
