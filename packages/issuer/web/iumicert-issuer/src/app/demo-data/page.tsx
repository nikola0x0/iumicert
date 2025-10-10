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
  Database,
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
                View Output
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
    <div className="space-y-8 pb-12">
      {/* Header with gradient background and patterns */}
      <div className="relative overflow-hidden bg-gradient-to-br from-blue-500 via-indigo-600 to-purple-600 rounded-3xl p-10 shadow-2xl shadow-blue-500/20">
        {/* Decorative elements */}
        <div className="absolute top-0 right-0 w-72 h-72 bg-white/5 rounded-full blur-3xl"></div>
        <div className="absolute -bottom-12 -left-12 w-64 h-64 bg-white/5 rounded-full blur-3xl"></div>
        <div className="absolute top-1/2 right-1/4 w-32 h-32 border-4 border-white/10 rounded-full"></div>
        <div className="absolute bottom-1/4 left-1/3 w-20 h-20 border-4 border-white/10 rounded-full"></div>

        <div className="relative z-10">
          <div className="inline-flex items-center gap-2 px-4 py-2 bg-white/10 backdrop-blur-sm rounded-full border border-white/20 mb-4">
            <Database className="w-4 h-4 text-white" />
            <span className="text-xs font-semibold text-white/90">DATA MANAGEMENT</span>
          </div>
          <h1 className="text-5xl font-extrabold text-white mb-3 tracking-tight">
            Demo Data Management
          </h1>
          <p className="text-blue-100 text-lg max-w-2xl leading-relaxed">
            Generate and manage test data for the IU-MiCert credential system
          </p>
        </div>
      </div>

      {/* System Reset Section */}
      <Card className="relative overflow-hidden rounded-3xl shadow-xl border-2 border-red-200 hover:border-red-300 transition-all duration-300 bg-white group">
        <div className="absolute top-0 left-0 w-2 h-full bg-gradient-to-b from-red-500 to-rose-600"></div>

        <CardHeader className="relative">
          <div className="flex items-start gap-4">
            <div className="relative">
              <div className="w-14 h-14 rounded-2xl bg-gradient-to-br from-red-500 to-rose-600 flex items-center justify-center shadow-lg shadow-red-500/30 transform group-hover:scale-110 transition-transform duration-300">
                <Trash2 className="w-7 h-7 text-white" />
              </div>
              <div className="absolute -top-1 -right-1 w-4 h-4 bg-red-500 rounded-full animate-ping"></div>
              <div className="absolute -top-1 -right-1 w-4 h-4 bg-red-500 rounded-full"></div>
            </div>
            <div className="flex-1">
              <div className="flex items-center gap-2 mb-2">
                <CardTitle className="text-2xl text-red-700">System Reset</CardTitle>
                <Badge variant="destructive" className="text-[10px] px-2 py-0.5">DESTRUCTIVE</Badge>
              </div>
              <CardDescription className="text-base text-gray-600">
                This action will permanently delete <strong className="text-red-600">all generated data</strong>, including terms,
                receipts, and Verkle trees. This operation cannot be undone.
              </CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-4 relative">
          <Button
            onClick={handleReset}
            disabled={resetStatus.status === "running"}
            variant="destructive"
            className="rounded-2xl shadow-lg hover:shadow-xl hover:shadow-red-500/30 transition-all transform hover:scale-[1.02] active:scale-[0.98]"
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
      <Card className="relative overflow-hidden rounded-3xl shadow-2xl border-2 border-blue-200 hover:border-blue-300 transition-all duration-300 bg-white group">

        <CardHeader className="relative">
          <div className="flex items-start gap-4">
            <div className="relative">
              <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-blue-500 via-indigo-600 to-purple-600 flex items-center justify-center shadow-xl shadow-blue-500/40 transform group-hover:rotate-6 transition-transform duration-300">
                <Zap className="w-8 h-8 text-white" />
              </div>
              <div className="absolute -bottom-2 -right-2 w-6 h-6 bg-yellow-400 rounded-full flex items-center justify-center shadow-md">
                <span className="text-xs">⚡</span>
              </div>
            </div>
            <div className="flex-1">
              <div className="flex items-center gap-3 mb-2">
                <CardTitle className="text-3xl bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent">
                  Data Generation
                </CardTitle>
              </div>
              <CardDescription className="text-base text-gray-600 leading-relaxed">
                Generate comprehensive test datasets with customizable parameters. Creates student journeys,
                converts to Verkle format, builds cryptographic trees, and generates verification receipts.
              </CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-8 relative">
          {/* Number of Students Input - Enhanced */}
          <div className="relative">
            <label className="flex items-center gap-2 text-sm font-semibold text-slate-800 mb-3">
              <span className="flex items-center justify-center w-6 h-6 rounded-full bg-blue-100 text-blue-600 text-xs">1</span>
              Number of Students
              <Badge variant="destructive" className="text-[10px] px-2 py-0.5">
                Required
              </Badge>
            </label>
            <div className="relative group">
              <input
                type="number"
                value={numStudentsForFull}
                onChange={(e) =>
                  setNumStudentsForFull(parseInt(e.target.value) || 1)
                }
                min="1"
                max="100"
                className="w-full px-5 py-4 border-2 border-gray-200 rounded-2xl focus:ring-4 focus:ring-blue-500/20 focus:border-blue-500 hover:border-blue-300 transition-all bg-white shadow-sm text-lg font-semibold text-gray-900 group-hover:shadow-md"
              />
              <div className="absolute right-4 top-1/2 -translate-y-1/2 text-gray-400 text-sm font-medium pointer-events-none">
                students
              </div>
            </div>
            <p className="text-sm text-gray-500 mt-3 ml-8">
              Range: 1-100 students, 3-6 courses each per term
            </p>
          </div>

          {/* Term Selection - Interactive Grid */}
          <div className="relative">
            <div className="flex items-center justify-between mb-4">
              <label className="flex items-center gap-2 text-sm font-semibold text-slate-800">
                <span className="flex items-center justify-center w-6 h-6 rounded-full bg-blue-100 text-blue-600 text-xs">2</span>
                Select Academic Terms
                <Badge variant="destructive" className="text-[10px] px-2 py-0.5">
                  Required
                </Badge>
              </label>
              <Button
                onClick={toggleAllTerms}
                variant="ghost"
                size="sm"
                className="rounded-xl hover:bg-blue-100 hover:text-blue-700 transition-colors text-xs font-semibold"
              >
                {selectedTerms.length === availableTerms.length
                  ? "✕ Deselect All"
                  : "✓ Select All"}
              </Button>
            </div>

            <div className="relative p-6 bg-gradient-to-br from-gray-50 via-blue-50/50 to-indigo-50/50 rounded-2xl border-2 border-gray-200 shadow-inner">
              <div className="grid grid-cols-2 gap-3">
                {availableTerms.map((term, idx) => (
                  <label
                    key={term}
                    className="relative flex items-center space-x-3 cursor-pointer p-4 rounded-xl transition-all duration-200 group bg-white border-2 border-gray-200 hover:border-blue-400 hover:shadow-lg hover:-translate-y-0.5"
                    style={{ animationDelay: `${idx * 50}ms` }}
                  >
                    <div className={`flex items-center justify-center w-5 h-5 rounded-md border-2 transition-all ${
                      selectedTerms.includes(term)
                        ? 'bg-gradient-to-br from-blue-500 to-indigo-600 border-blue-500'
                        : 'border-gray-300 group-hover:border-blue-400'
                    }`}>
                      {selectedTerms.includes(term) && (
                        <svg className="w-3 h-3 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
                        </svg>
                      )}
                    </div>
                    <input
                      type="checkbox"
                      checked={selectedTerms.includes(term)}
                      onChange={() => toggleTerm(term)}
                      className="sr-only"
                    />
                    <span className={`text-sm font-semibold transition-colors ${
                      selectedTerms.includes(term)
                        ? 'text-blue-700'
                        : 'text-slate-700 group-hover:text-blue-600'
                    }`}>
                      {term.replace(/_/g, ' ')}
                    </span>
                  </label>
                ))}
              </div>
            </div>

            <div className="mt-3 ml-8 text-sm">
              <div className={`inline-flex px-3 py-1.5 rounded-full font-semibold ${
                selectedTerms.length === 0
                  ? 'bg-red-100 text-red-700'
                  : 'bg-green-100 text-green-700'
              }`}>
                {selectedTerms.length} / {availableTerms.length} terms selected
              </div>
            </div>
          </div>

          {/* Action Button - Hero CTA */}
          <div className="pt-4">
            <Button
              onClick={handleGenerateFull}
              disabled={
                generateStatus.status === "running" || selectedTerms.length === 0
              }
              className="w-full text-white rounded-2xl bg-gradient-to-r from-blue-600 via-indigo-600 to-purple-600 hover:from-blue-700 hover:via-indigo-700 hover:to-purple-700 shadow-2xl shadow-blue-500/50 hover:shadow-blue-600/60 transition-all transform hover:scale-[1.02] active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none group relative overflow-hidden"
              size="lg"
            >
              <div className="absolute inset-0 bg-gradient-to-r from-white/0 via-white/20 to-white/0 translate-x-[-100%] group-hover:translate-x-[100%] transition-transform duration-1000"></div>
              <div className="relative flex items-center justify-center gap-3 py-1">
                {generateStatus.status === "running" ? (
                  <>
                    <Loader2 className="w-6 h-6 animate-spin" />
                    <span className="text-lg font-bold">Generating...</span>
                  </>
                ) : (
                  <>
                    <Zap className="w-6 h-6 group-hover:rotate-12 transition-transform" />
                    <span className="text-lg font-bold">
                      Generate {numStudentsForFull} Students × {selectedTerms.length} Terms
                    </span>
                    <div className="absolute right-6 top-1/2 -translate-y-1/2 opacity-0 group-hover:opacity-100 transition-opacity">
                      →
                    </div>
                  </>
                )}
              </div>
            </Button>
          </div>

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
