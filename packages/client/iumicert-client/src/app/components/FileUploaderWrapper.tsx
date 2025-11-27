"use client";

import React from "react";
import { FileUploader } from "react-drag-drop-files";
import { StyleSheetManager } from "styled-components";
import isPropValid from "@emotion/is-prop-valid";

// Custom prop validation to filter out problematic props
const shouldForwardProp = (prop: string) => {
  // Don't forward these problematic props
  if (["override"].includes(prop)) {
    return false;
  }
  // Use emotion's default validation for other props
  return isPropValid(prop);
};

interface FileUploaderWrapperProps {
  handleChange: (file: File) => void;
  name: string;
  types: string[];
  onTypeError?: (error: string) => void;
  maxSize?: number;
  classes?: string;
  children: React.ReactNode;
}

const FileUploaderWrapper: React.FC<FileUploaderWrapperProps> = (props) => {
  return (
    <StyleSheetManager shouldForwardProp={shouldForwardProp}>
      <FileUploader {...props} />
    </StyleSheetManager>
  );
};

export default FileUploaderWrapper;
export { FileUploaderWrapper };
