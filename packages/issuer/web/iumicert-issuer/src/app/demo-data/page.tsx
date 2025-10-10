"use client";

import { useState } from "react";
import {
  Trash2,
  Zap,
  PlusCircle,
  Loader2,
  CheckCircle2,
  XCircle,
  AlertCircle,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";

type OperationStatus = {
  status: "idle" | "running" | "success" | "error";
  message?: string;
  output?: string;
};

export default function DemoDataPage() {
  const [resetStatus, setResetStatus] = useState<OperationStatus>({
    status: "idle",
  });
  const [generateStatus, setGenerateStatus] = useState<OperationStatus>({
    status: "idle",
  });
  const [addTermStatus, setAddTermStatus] = useState<OperationStatus>({
    status: "idle",
  });

  // Add New Term state
  const [termId, setTermId] = useState("Semester_1_2026");

  // Full Data Generation state
  const [numStudentsForFull, setNumStudentsForFull] = useState(5);
  const [availableTerms, setAvailableTerms] = useState<string[]>([
    "Semester_1_2023",
    "Semester_2_2023",
    "Summer_2023",
    "Semester_1_2024",
    "Semester_2_2024",
    "Summer_2024",
    "Semester_1_2025",
  ]);
  const [selectedTerms, setSelectedTerms] = useState<string[]>([
    "Semester_1_2023",
    "Semester_2_2023",
    "Summer_2023",
    "Semester_1_2024",
    "Semester_2_2024",
    "Summer_2024",
    "Semester_1_2025",
  ]);

  const handleReset = async () => {
    if (
      !confirm("⚠️ WARNING: This will delete ALL generated data. Are you sure?")
    ) {
      return;
    }

    setResetStatus({ status: "running", message: "Executing reset.sh..." });

    try {
      const response = await fetch("http://localhost:8080/api/demo/reset", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
      });

      const data = await response.json();

      if (!response.ok || !data.success) {
        throw new Error(data.error || "Reset failed");
      }

      setResetStatus({
        status: "success",
        message: "System reset completed successfully",
        output: data.data?.output,
      });
    } catch (error: any) {
      setResetStatus({
        status: "error",
        message: error.message || "Failed to execute reset",
      });
    }
  };

  const handleGenerateFull = async () => {
    if (selectedTerms.length === 0) {
      alert("Please select at least one term to generate");
      return;
    }

    setGenerateStatus({
      status: "running",
      message: `Generating ${numStudentsForFull} students across ${selectedTerms.length} terms...`,
    });

    try {
      const response = await fetch(
        "http://localhost:8080/api/demo/generate-full",
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            num_students: numStudentsForFull,
            terms: selectedTerms,
          }),
        }
      );

      const data = await response.json();

      if (!response.ok || !data.success) {
        throw new Error(data.error || "Generation failed");
      }

      setGenerateStatus({
        status: "success",
        message: data.data?.message || "Full dataset generated successfully",
        output: data.data?.output,
      });
    } catch (error: any) {
      setGenerateStatus({
        status: "error",
        message: error.message || "Failed to execute generation",
      });
    }
  };

  const toggleTerm = (term: string) => {
    setSelectedTerms((prev) =>
      prev.includes(term) ? prev.filter((t) => t !== term) : [...prev, term]
    );
  };

  const toggleAllTerms = () => {
    setSelectedTerms((prev) =>
      prev.length === availableTerms.length ? [] : [...availableTerms]
    );
  };

  const handleAddTerm = () => {
    if (!termId.trim()) {
      alert("Please enter a term ID");
      return;
    }

    if (availableTerms.includes(termId)) {
      alert("This term already exists in the list");
      return;
    }

    // Add term to available terms and automatically select it
    setAvailableTerms((prev) => [...prev, termId]);
    setSelectedTerms((prev) => [...prev, termId]);

    // Clear the input and show success
    setTermId("");
    setAddTermStatus({
      status: "success",
      message: `Added "${termId}" to term list and selected it`,
    });

    // Clear success message after 3 seconds
    setTimeout(() => {
      setAddTermStatus({ status: "idle" });
    }, 3000);
  };

  const StatusBadge = ({ status, message, output }: OperationStatus) => {
    if (status === "idle") return null;

    const config = {
      running: {
        variant: "default" as const,
        icon: Loader2,
        className:
          "bg-gradient-to-r from-blue-50 to-indigo-50 border-blue-200 text-blue-900 shadow-sm",
        iconBg: "bg-blue-100",
        iconColor: "text-blue-600",
      },
      success: {
        variant: "default" as const,
        icon: CheckCircle2,
        className:
          "bg-gradient-to-r from-green-50 to-emerald-50 border-green-200 text-green-900 shadow-sm",
        iconBg: "bg-green-100",
        iconColor: "text-green-600",
      },
      error: {
        variant: "destructive" as const,
        icon: XCircle,
        className:
          "bg-gradient-to-r from-red-50 to-rose-50 border-red-200 text-red-900 shadow-sm",
        iconBg: "bg-red-100",
        iconColor: "text-red-600",
      },
    };

    const { icon: Icon, className, iconBg, iconColor } = config[status];

    return (
      <Alert className={`${className} mt-4 rounded-xl flex items-start gap-3`}>
        <div
          className={`w-10 h-10 rounded-lg ${iconBg} flex items-center justify-center flex-shrink-0`}
        >
          <Icon
            className={`h-5 w-5 ${iconColor} ${
              status === "running" ? "animate-spin" : ""
            }`}
          />
        </div>
        <AlertDescription className="flex-1">
          <p className="font-semibold">{message}</p>
          {output && (
            <details className="mt-3">
              <summary className="cursor-pointer text-sm font-medium hover:text-blue-600 transition-colors">
                📄 View Output
              </summary>
              <pre className="mt-2 p-3 bg-white rounded-lg text-xs overflow-x-auto border shadow-sm">
                {output}
              </pre>
            </details>
          )}
        </AlertDescription>
      </Alert>
    );
  };

  return (
    <div className="max-w-6xl mx-auto space-y-8">
      {/* Header with gradient background */}
      <div className="relative overflow-hidden bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 rounded-2xl p-8 border border-blue-100 shadow-sm">
        <div className="absolute top-0 right-0 w-64 h-64 bg-blue-400/10 rounded-full blur-3xl -z-10"></div>
        <div className="absolute bottom-0 left-0 w-48 h-48 bg-indigo-400/10 rounded-full blur-3xl -z-10"></div>
        <h1 className="text-4xl font-bold bg-gradient-to-r from-blue-600 via-indigo-600 to-purple-600 bg-clip-text text-transparent">
          Demo Data Management
        </h1>
        <p className="text-gray-700 mt-3 text-lg">
          Generate and manage demo data for testing the IU-MiCert system
        </p>
      </div>

      {/* System Reset Section */}
      <Card className="border-l-4 border-red-500 rounded-2xl shadow-lg hover:shadow-xl transition-shadow duration-300 bg-gradient-to-br from-white to-red-50/30">
        <CardHeader>
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 rounded-xl bg-red-100 flex items-center justify-center">
              <Trash2 className="w-6 h-6 text-red-600" />
            </div>
            <div>
              <CardTitle className="text-red-700">System Reset</CardTitle>
              <CardDescription className="mt-1">
                <strong>Warning:</strong> This will delete ALL generated data,
                including terms, receipts, and Verkle trees.
              </CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          <Button
            onClick={handleReset}
            disabled={resetStatus.status === "running"}
            variant="destructive"
            className="rounded-xl shadow-md hover:shadow-lg hover:shadow-red-500/20 transition-all"
          >
            {resetStatus.status === "running" ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Resetting...
              </>
            ) : (
              <>
                <Trash2 className="mr-2 h-4 w-4" />
                Execute Reset (./reset.sh)
              </>
            )}
          </Button>

          <StatusBadge {...resetStatus} />
        </CardContent>
      </Card>

      {/* Full Data Generation Section */}
      <Card className="border-l-4 border-blue-500 rounded-2xl shadow-lg hover:shadow-xl transition-shadow duration-300 bg-gradient-to-br from-white via-blue-50/20 to-indigo-50/30">
        <CardHeader>
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center shadow-md">
              <Zap className="w-6 h-6 text-white" />
            </div>
            <div>
              <CardTitle className="text-blue-700">
                Full Data Generation
              </CardTitle>
              <CardDescription className="mt-1">
                Generate a complete dataset with custom parameters. This creates
                student journeys, converts them to Verkle format, builds trees,
                and generates receipts for all selected terms.
              </CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-6">
          {/* Number of Students Input */}
          <div>
            <label className="flex items-center gap-2 text-sm font-medium text-slate-700 mb-2">
              Number of Students
              <Badge variant="destructive" className="text-xs">
                Required
              </Badge>
            </label>
            <input
              type="number"
              value={numStudentsForFull}
              onChange={(e) =>
                setNumStudentsForFull(parseInt(e.target.value) || 1)
              }
              min="1"
              max="100"
              className="w-full px-4 py-3 border-2 border-gray-200 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 hover:border-gray-300 transition-all bg-white shadow-sm"
            />
            <p className="text-xs text-gray-500 mt-2">
              Number of students to generate (1-100). Each student gets 3-6
              courses per term.
            </p>
          </div>

          {/* Term Selection */}
          <div>
            <div className="flex items-center justify-between mb-3">
              <div className="flex items-center gap-2">
                <label className="text-sm font-medium text-slate-700">
                  Select Terms
                </label>
                <Badge variant="destructive" className="text-xs">
                  Required
                </Badge>
              </div>
              <Button
                onClick={toggleAllTerms}
                variant="ghost"
                size="sm"
                className="rounded-lg"
              >
                {selectedTerms.length === availableTerms.length
                  ? "Deselect All"
                  : "Select All"}
              </Button>
            </div>
            <div className="grid grid-cols-2 gap-3 p-5 bg-gradient-to-br from-slate-50 to-blue-50/30 rounded-xl border border-blue-100 shadow-sm">
              {availableTerms.map((term) => (
                <label
                  key={term}
                  className="flex items-center space-x-3 cursor-pointer hover:bg-white p-3 rounded-lg transition-all duration-200 hover:shadow-md group"
                >
                  <input
                    type="checkbox"
                    checked={selectedTerms.includes(term)}
                    onChange={() => toggleTerm(term)}
                    className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500 cursor-pointer"
                  />
                  <span className="text-sm font-medium text-slate-900 group-hover:text-blue-600 transition-colors">
                    {term}
                  </span>
                </label>
              ))}
            </div>
            <p className="text-xs text-gray-500 mt-2">
              Selected: {selectedTerms.length} / {availableTerms.length} terms
            </p>
          </div>

          {/* Action Button */}
          <Button
            onClick={handleGenerateFull}
            disabled={
              generateStatus.status === "running" || selectedTerms.length === 0
            }
            className="w-full text-white rounded-xl bg-gradient-to-r from-blue-600 via-indigo-600 to-blue-600 hover:from-blue-700 hover:via-indigo-700 hover:to-blue-700 shadow-lg shadow-blue-500/30 hover:shadow-xl hover:shadow-blue-500/40 transition-all"
            size="lg"
          >
            {generateStatus.status === "running" ? (
              <>
                <Loader2 className="mr-2 h-5 w-5 animate-spin" />
                Generating...
              </>
            ) : (
              <>
                <Zap className="mr-2 h-5 w-5" />
                Generate {numStudentsForFull} Students × {selectedTerms.length}{" "}
                Terms
              </>
            )}
          </Button>

          <StatusBadge {...generateStatus} />

          {/* Add New Term to List */}
          <div className="pt-6 border-t-2 border-dashed border-gray-200 mt-8">
            <div className="flex items-center gap-3 mb-4">
              <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-green-500 to-emerald-600 flex items-center justify-center shadow-md">
                <PlusCircle className="w-5 h-5 text-white" />
              </div>
              <div>
                <h3 className="text-lg font-bold text-slate-900">
                  Add Custom Term
                </h3>
                <p className="text-gray-600 text-sm">
                  Add a new term to the selection list above
                </p>
              </div>
            </div>

            <div className="space-y-4">
              <div>
                <label className="flex items-center gap-2 text-sm font-medium text-slate-700 mb-2">
                  New Term ID
                  <Badge variant="destructive" className="text-xs">
                    Required
                  </Badge>
                </label>
                <input
                  type="text"
                  value={termId}
                  onChange={(e) => setTermId(e.target.value)}
                  onKeyDown={(e) => {
                    if (e.key === "Enter") {
                      handleAddTerm();
                    }
                  }}
                  placeholder="e.g., Semester_1_2026"
                  className="w-full px-4 py-3 border-2 border-gray-200 rounded-xl focus:ring-2 focus:ring-green-500 focus:border-green-500 hover:border-gray-300 transition-all bg-white shadow-sm"
                />
                <p className="text-xs text-gray-500 mt-2">
                  Format: Semester_1_YYYY, Semester_2_YYYY, or Summer_YYYY
                </p>
              </div>

              <Button
                onClick={handleAddTerm}
                disabled={!termId.trim()}
                variant="default"
                className="bg-gradient-to-r from-green-600 to-emerald-600 hover:from-green-700 hover:to-emerald-700 text-white rounded-xl w-full shadow-lg shadow-green-500/30 hover:shadow-xl hover:shadow-green-500/40 transition-all"
                size="lg"
              >
                <PlusCircle className="mr-2 h-5 w-5" />
                Add Term
              </Button>
            </div>

            <StatusBadge {...addTermStatus} />
          </div>
        </CardContent>
      </Card>

      {/* Info Section */}
      <Alert className="bg-gradient-to-r from-blue-50 to-indigo-50 border-blue-200 rounded-xl shadow-sm">
        <div className="flex items-center">
          <AlertCircle className="h-5 w-5 text-blue-600" />
          <AlertDescription className="ml-2 text-blue-900">
            After generating data, use the main issuer dashboard to process
            terms, build Verkle trees, and publish to blockchain.
          </AlertDescription>
        </div>
      </Alert>
    </div>
  );
}
