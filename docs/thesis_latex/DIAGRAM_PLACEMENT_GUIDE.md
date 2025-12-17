# Draw.io Diagram Placement Guide

## Export Instructions
1. Open each `.drawio` file in https://app.diagrams.net
2. File → Export As → PNG
3. Settings: Scale 200%, White Background, Border Width 5
4. Save to `figures/` directory with the name shown below

## Diagram Mapping

| Draw.io File | Export As | Use In | Replaces |
|--------------|-----------|--------|----------|
| `verkle_microcred_structure_updated.drawio` | `verkle_microcred_structure.png` | **method.tex** (optional - can add after line 277) | Old hierarchical structure |
| `01_system_architecture.drawio` | `system_architecture_3tier.png` | **method.tex** line 135 (already using system_overview.png) | Can use as alternative |
| `02_verkle_tree_structure.drawio` | `verkle_tree_academic.png` | **related.tex** line 109 (already using verkle_tree_structure.png) | Can use as alternative |
| `03_data_pipeline.drawio` | `data_pipeline.png` | **prototyping.tex** line 223 | ASCII verbatim block |
| `04_revocation_versioning.drawio` | `revocation_versioning.png` | **method.tex** line 356 | ASCII verbatim block |
| `05_component_interaction.drawio` | `component_interaction.png` | **method.tex** line 151 (already using system_architecture.png) | Can use as alternative |

## Recommended Actions

### MUST EXPORT (to replace ASCII):
1. **`04_revocation_versioning.drawio`** → Export as `revocation_versioning.png`
   - Replace ASCII verbatim at method.tex:356

2. **`03_data_pipeline.drawio`** → Export as `data_pipeline.png`
   - Replace ASCII verbatim at prototyping.tex:223

### OPTIONAL (alternatives to existing PNGs):
The other diagrams are alternatives to existing PNG figures that are already referenced.

### NEW DIAGRAM (recommended to add):
3. **`verkle_microcred_structure_updated.drawio`** → Export as `verkle_microcred_structure.png`
   - Add to method.tex after line 277 (in "Tree Structure" section)
   - This shows the correct "one tree per term" design

## LaTeX Update Examples

### For revocation diagram (method.tex line 354):
```latex
\begin{figure}[H]
\centering
\includegraphics[width=0.9\textwidth]{figures/revocation_versioning.png}
\caption{Version-Based Revocation Through Tree Supersession}
\label{fig:revocation-concept}
\end{figure}
```

### For data pipeline (prototyping.tex line 221):
```latex
\begin{figure}[H]
\centering
\includegraphics[width=0.85\textwidth]{figures/data_pipeline.png}
\caption{IU-MiCert Data Pipeline}
\label{fig:data-pipeline}
\end{figure}
```

### For verkle microcred structure (method.tex, add after line 277):
```latex
\begin{figure}[H]
\centering
\includegraphics[width=0.9\textwidth]{figures/verkle_microcred_structure.png}
\caption{One Verkle Tree Per Academic Term}
\label{fig:verkle-microcred-structure}
\end{figure}
```
