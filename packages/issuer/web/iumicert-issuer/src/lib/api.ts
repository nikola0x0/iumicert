const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export interface Term {
  id: string;
  name: string;
  start_date: string;
  end_date: string;
  status: string;
  student_count?: number;
  total_courses?: number;
}

export interface Receipt {
  id: string;
  student_id: string;
  student_name: string;
  term_id: string;
  courses: Course[];
  merkle_root: string;
  verkle_proof: string;
  created_at: string;
}

export interface Course {
  id: string;
  name: string;
  grade: string;
  credits: number;
}

export interface BlockchainHistory {
  transaction_hash: string;
  block_number: number;
  timestamp: string;
  operation: string;
  data: any;
}

interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}

class ApiService {
  private async request<T>(endpoint: string, options?: RequestInit): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`;
    const response = await fetch(url, {
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
      ...options,
    });

    if (!response.ok) {
      throw new Error(`API Error: ${response.status} ${response.statusText}`);
    }

    const result: ApiResponse<T> = await response.json();
    
    if (!result.success) {
      throw new Error(result.error || 'API request failed');
    }
    
    return result.data as T;
  }

  // Term management
  async getTerms(): Promise<Term[]> {
    return this.request<Term[]>('/api/terms');
  }

  async createTerm(term: Omit<Term, 'id'>): Promise<Term> {
    return this.request<Term>('/api/terms', {
      method: 'POST',
      body: JSON.stringify(term),
    });
  }

  async getTermRoot(termId: string): Promise<{ verkle_root: string; merkle_root: string; total_students: number }> {
    return this.request<{ verkle_root: string; merkle_root: string; total_students: number }>(`/api/terms/${termId}/roots`);
  }

  // Receipt management
  async getReceipts(termId: string): Promise<Receipt[]> {
    return this.request<Receipt[]>(`/api/terms/${termId}/receipts`);
  }

  async generateReceipt(data: {
    student_id: string;
    student_name: string;
    term_id: string;
    courses: Course[];
  }): Promise<Receipt> {
    return this.request<Receipt>('/api/receipts', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  // Verification
  async verifyLocal(data: {
    receipt_id: string;
    merkle_proof: string;
    verkle_proof: string;
  }): Promise<{ valid: boolean; details: string }> {
    return this.request<{ valid: boolean; details: string }>('/api/verify-local', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  // Blockchain operations
  async publishRoots(data: {
    merkle_root: string;
    verkle_root: string;
    receipts: string[];
  }): Promise<{ transaction_hash: string }> {
    return this.request<{ transaction_hash: string }>('/api/publish-roots', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async publishToBlockchain(data: {
    term_id: string;
    network: string;
    gas_limit: number;
  }): Promise<{ transaction_hash: string; status: string }> {
    return this.request<{ transaction_hash: string; status: string }>('/api/blockchain/publish', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async getBlockchainHistory(): Promise<BlockchainHistory[]> {
    return this.request<BlockchainHistory[]>('/api/blockchain/history');
  }

  // Health check
  async healthCheck(): Promise<{ status: string; timestamp: string }> {
    try {
      const url = `${API_BASE_URL}/api/health`;
      const response = await fetch(url);
      
      if (!response.ok) {
        return { status: 'error', timestamp: new Date().toISOString() };
      }
      
      // Health endpoint might return plain text or different format
      const result = await response.json();
      return result.data || result;
    } catch (error) {
      return { status: 'error', timestamp: new Date().toISOString() };
    }
  }
}

export const apiService = new ApiService();