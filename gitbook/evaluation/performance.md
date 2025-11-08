# Performance Evaluation

## Key Metrics

IU-MiCert demonstrates practical performance for real-world academic credential management.

## Proof Size Comparison

| System | Proof Size per Course | Student with 120 Courses |
|--------|----------------------|--------------------------|
| Merkle Tree | ~224 bytes | ~27 KB total |
| **IU-MiCert (Verkle)** | **~32 bytes** | **~4 KB total** |
| **Reduction** | **86%** | **85%** |

**Significance**: Enables mobile verification, reduces bandwidth costs

## Verification Performance

Measured on deployed system:

- **Proof verification**: <100ms
- **Blockchain query**: <500ms
- **Total verification**: <1 second
- **API response**: ~200ms average

**Comparison**: Traditional systems require contacting institution (hours to days)

## Storage Requirements

### On-Chain (Blockchain)

- **Per term**: 32 bytes (one root commitment)
- **6 terms**: 192 bytes total
- **Gas cost**: ~$0.10 USD per term (at 20 gwei)

**Efficiency**: Constant storage regardless of student count

### Off-Chain (Receipts)

- **Per receipt**: ~2-5 KB JSON
- **Includes**: Course details + 32-byte proof + metadata
- **Storage**: Student-managed (USB drive, cloud, email)

## Scalability Analysis

### Tree Construction Time

For 5,000 students × 6 courses per term = 30,000 leaves:

- **Tree building**: ~5-10 seconds
- **Proof generation**: ~1ms per course
- **All receipts**: ~30 seconds total

**Feasibility**: Practical for semester-end batch processing

### Verification Throughput

- **Sequential**: ~10 verifications/second
- **Parallel**: Limited by blockchain RPC calls
- **Caching**: Repeated verifications <10ms (cached roots)

**Conclusion**: Supports thousands of verifications per day

## Bandwidth Requirements

| Operation | Data Transfer |
|-----------|--------------|
| Upload receipt | ~2-5 KB |
| Blockchain query | ~1 KB |
| Verification response | ~0.5 KB |
| **Total** | **~3-7 KB** |

**Mobile-friendly**: Works on 3G networks

## Comparison with Related Work

| System | Proof Size | Verification Time | Blockchain Cost |
|--------|-----------|-------------------|-----------------|
| Blockcerts | ~2 KB | <1s | $0.50/credential |
| **IU-MiCert** | **32 bytes** | **<1s** | **$0.10/term** |
| Traditional | N/A | Hours-days | N/A |

**Advantage**: Lower cost, smaller proofs, same speed

## Resource Utilization

### Issuer System

- **CPU**: Low (batch processing)
- **Memory**: ~2 GB during tree construction
- **Storage**: ~100 MB per 10,000 students
- **Network**: Minimal (one blockchain tx per term)

### Verification Portal

- **CPU**: Minimal (cryptographic ops in <100ms)
- **Memory**: <50 MB
- **Network**: 3-7 KB per verification

**Deployment**: Runs on standard cloud infrastructure

## Real-World Performance

From deployed system on Ethereum Sepolia:

- **Terms published**: 6 (Semester_1_2023 through Summer_2024)
- **Total gas used**: ~300,000 gas
- **Cost**: ~$0.60 USD total (at 20 gwei)
- **Receipts generated**: 30+ test students
- **Verification uptime**: 99%+ (Vercel hosting)

## Bottlenecks & Limitations

### Current Bottlenecks

1. **Blockchain RPC**: Limited by provider rate limits
2. **Proof generation**: Sequential processing per student
3. **Database queries**: Large term retrievals

### Optimization Opportunities

1. **Caching**: Store frequently accessed roots locally
2. **Parallel processing**: Generate proofs concurrently
3. **Batch blockchain queries**: Fetch multiple roots at once
4. **CDN**: Serve static verification portal globally

## Future Performance Enhancements

- **Layer 2 deployment**: Reduce gas costs 10-100x
- **Proof batching**: Combine multiple course proofs
- **Incremental trees**: Update trees without full rebuild
- **WASM verification**: Browser-based proof checking

---

**Next**: Comparison with existing credential systems.
