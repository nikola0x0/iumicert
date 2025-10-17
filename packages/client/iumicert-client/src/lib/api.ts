const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}

class ApiService {
  private async request<T>(
    endpoint: string,
    options?: RequestInit
  ): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`;
    const response = await fetch(url, {
      headers: {
        "Content-Type": "application/json",
        ...options?.headers,
      },
      ...options,
    });

    if (!response.ok) {
      throw new Error(`API Error: ${response.status} ${response.statusText}`);
    }

    const result: ApiResponse<T> = await response.json();

    if (!result.success) {
      throw new Error(result.error || "API request failed");
    }

    return result.data as T;
  }

  // Full IPA verification
  async verifyReceiptIPA(receipt: any): Promise<{
    status: string;
    student_id: string;
    total_courses: number;
    verified_courses: number;
    failed_courses: number;
    failed_list: string[];
    term_results: Record<string, any>;
    computation_note: string;
  }> {
    return this.request<{
      status: string;
      student_id: string;
      total_courses: number;
      verified_courses: number;
      failed_courses: number;
      failed_list: string[];
      term_results: Record<string, any>;
      computation_note: string;
    }>("/api/verifier/ipa-verify", {
      method: "POST",
      body: JSON.stringify({ receipt }),
    });
  }

  // Health check
  async healthCheck(): Promise<{ status: string; timestamp: string }> {
    try {
      const url = `${API_BASE_URL}/api/health`;
      const response = await fetch(url);

      if (!response.ok) {
        return { status: "error", timestamp: new Date().toISOString() };
      }

      const result = await response.json();
      return result.data || result;
    } catch (error) {
      return { status: "error", timestamp: new Date().toISOString() };
    }
  }
}

export const apiService = new ApiService();
