import React, { useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  Alert,
} from 'react-native';
import { useLocalSearchParams, router } from 'expo-router';
import * as DocumentPicker from 'expo-document-picker';
import { blueprintsApi } from '../../../../../src/api/blueprints';
import { Button } from '../../../../../src/components/ui/Button';
import { Card } from '../../../../../src/components/ui/Card';
import { COLORS, MAX_FILE_SIZE, SUPPORTED_FILE_TYPES } from '../../../../../src/utils/constants';
import { useRequestUploadUrl, useCompleteUpload } from '../../../../../src/hooks/useBlueprints';

type UploadStep = 'idle' | 'selecting' | 'getting-url' | 'uploading' | 'finalizing' | 'complete' | 'error';

interface SelectedFile {
  name: string;
  size: number;
  type: string;
  uri: string;
}

export default function BlueprintUploadScreen() {
  const { id: projectId } = useLocalSearchParams<{ id: string }>();
  const [step, setStep] = useState<UploadStep>('idle');
  const [selectedFile, setSelectedFile] = useState<SelectedFile | null>(null);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [blueprintId, setBlueprintId] = useState<string | null>(null);

  const requestUploadUrl = useRequestUploadUrl();
  const completeUpload = useCompleteUpload();

  const formatFileSize = (bytes: number): string => {
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(2)} KB`;
    return `${(bytes / (1024 * 1024)).toFixed(2)} MB`;
  };

  const validateFile = (file: DocumentPicker.DocumentPickerAsset): string | null => {
    // Check file size
    if (file.size && file.size > MAX_FILE_SIZE) {
      return `File size exceeds maximum allowed size of ${formatFileSize(MAX_FILE_SIZE)}`;
    }

    // Check file type
    const supportedTypes = Object.keys(SUPPORTED_FILE_TYPES);
    if (file.mimeType && !supportedTypes.includes(file.mimeType)) {
      return 'Unsupported file type. Please upload PDF, PNG, or JPG files.';
    }

    return null;
  };

  const handleFilePick = async () => {
    try {
      setStep('selecting');
      setError(null);

      const result = await DocumentPicker.getDocumentAsync({
        type: ['application/pdf', 'image/png', 'image/jpeg'],
        copyToCacheDirectory: true,
      });

      if (result.canceled) {
        setStep('idle');
        return;
      }

      const file = result.assets[0];
      const validationError = validateFile(file);

      if (validationError) {
        setError(validationError);
        setStep('error');
        return;
      }

      setSelectedFile({
        name: file.name,
        size: file.size || 0,
        type: file.mimeType || 'application/octet-stream',
        uri: file.uri,
      });
      setStep('idle');
    } catch (err) {
      console.error('Error picking file:', err);
      setError('Failed to select file. Please try again.');
      setStep('error');
    }
  };

  const handleUpload = async () => {
    if (!selectedFile || !projectId) {
      return;
    }

    try {
      // Step 1: Get upload URL
      setStep('getting-url');
      setUploadProgress(10);

      const uploadUrlData = await requestUploadUrl.mutateAsync({
        projectId,
        data: {
          filename: selectedFile.name,
          content_type: selectedFile.type,
        },
      });

      setBlueprintId(uploadUrlData.blueprint_id);
      setUploadProgress(20);

      // Step 2: Upload to S3
      setStep('uploading');

      // Convert file URI to blob for upload
      const response = await fetch(selectedFile.uri);
      const blob = await response.blob();
      await blueprintsApi.uploadToS3(uploadUrlData.upload_url, blob);

      setUploadProgress(80);

      // Step 3: Complete upload
      setStep('finalizing');
      await completeUpload.mutateAsync({
        blueprintId: uploadUrlData.blueprint_id,
        success: true,
      });

      setUploadProgress(100);
      setStep('complete');

      // Show success message and navigate back
      setTimeout(() => {
        Alert.alert(
          'Success',
          'Blueprint uploaded successfully!',
          [
            {
              text: 'OK',
              onPress: () => router.back(),
            },
          ]
        );
      }, 500);
    } catch (err) {
      console.error('Upload error:', err);
      setError('Failed to upload blueprint. Please try again.');
      setStep('error');

      // Try to mark upload as failed in the backend
      if (blueprintId) {
        try {
          await completeUpload.mutateAsync({
            blueprintId,
            success: false,
            error_message: err instanceof Error ? err.message : 'Upload failed',
          });
        } catch (completeErr) {
          console.error('Failed to mark upload as failed:', completeErr);
        }
      }
    }
  };

  const handleReset = () => {
    setStep('idle');
    setSelectedFile(null);
    setUploadProgress(0);
    setError(null);
    setBlueprintId(null);
  };

  const getStepMessage = (): string => {
    switch (step) {
      case 'selecting':
        return 'Selecting file...';
      case 'getting-url':
        return 'Getting upload URL...';
      case 'uploading':
        return `Uploading... ${uploadProgress}%`;
      case 'finalizing':
        return 'Finalizing...';
      case 'complete':
        return 'Upload complete!';
      case 'error':
        return error || 'An error occurred';
      default:
        return 'Ready to upload';
    }
  };

  const isProcessing = ['selecting', 'getting-url', 'uploading', 'finalizing'].includes(step);

  return (
    <ScrollView style={styles.container}>
      <Card style={styles.card}>
        <Text style={styles.title}>Upload Blueprint</Text>
        <Text style={styles.subtitle}>
          Upload PDF, PNG, or JPG files (max 50MB)
        </Text>

        <View style={styles.section}>
          {!selectedFile ? (
            <Button
              title="Select File"
              onPress={handleFilePick}
              disabled={isProcessing}
            />
          ) : (
            <>
              <Card style={styles.fileCard}>
                <Text style={styles.fileName}>{selectedFile.name}</Text>
                <Text style={styles.fileInfo}>
                  Size: {formatFileSize(selectedFile.size)}
                </Text>
                <Text style={styles.fileInfo}>Type: {selectedFile.type}</Text>
              </Card>

              <Button
                title="Change File"
                onPress={handleFilePick}
                variant="secondary"
                disabled={isProcessing}
                style={styles.changeButton}
              />
            </>
          )}
        </View>

        {selectedFile && (
          <View style={styles.section}>
            <View style={styles.statusCard}>
              <Text style={styles.statusText}>{getStepMessage()}</Text>
              {isProcessing && (
                <View style={styles.progressBar}>
                  <View
                    style={[
                      styles.progressFill,
                      { width: `${uploadProgress}%` },
                    ]}
                  />
                </View>
              )}
            </View>

            {step === 'idle' && (
              <Button
                title="Upload"
                onPress={handleUpload}
                style={styles.uploadButton}
              />
            )}

            {step === 'complete' && (
              <Text style={styles.successText}>✓ Upload successful!</Text>
            )}

            {step === 'error' && (
              <>
                <Text style={styles.errorText}>{error}</Text>
                <Button
                  title="Retry"
                  onPress={handleReset}
                  style={styles.retryButton}
                />
              </>
            )}

            {isProcessing && (
              <Button
                title="Cancel"
                onPress={handleReset}
                variant="danger"
                style={styles.cancelButton}
              />
            )}
          </View>
        )}

        <Button
          title="Back to Project"
          onPress={() => router.back()}
          variant="secondary"
          disabled={isProcessing}
          style={styles.backButton}
        />
      </Card>
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.background.secondary,
  },
  card: {
    margin: 16,
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    color: COLORS.text.primary,
    marginBottom: 8,
  },
  subtitle: {
    fontSize: 14,
    color: COLORS.text.secondary,
    marginBottom: 24,
  },
  section: {
    marginBottom: 24,
  },
  fileCard: {
    backgroundColor: COLORS.background.secondary,
    marginBottom: 12,
  },
  fileName: {
    fontSize: 16,
    fontWeight: '600',
    color: COLORS.text.primary,
    marginBottom: 8,
  },
  fileInfo: {
    fontSize: 14,
    color: COLORS.text.secondary,
    marginBottom: 4,
  },
  changeButton: {
    marginTop: 8,
  },
  statusCard: {
    padding: 16,
    backgroundColor: COLORS.background.secondary,
    borderRadius: 8,
    marginBottom: 16,
  },
  statusText: {
    fontSize: 16,
    fontWeight: '600',
    color: COLORS.text.primary,
    textAlign: 'center',
    marginBottom: 12,
  },
  progressBar: {
    height: 8,
    backgroundColor: COLORS.border,
    borderRadius: 4,
    overflow: 'hidden',
  },
  progressFill: {
    height: '100%',
    backgroundColor: COLORS.primary,
    borderRadius: 4,
  },
  uploadButton: {
    marginTop: 8,
  },
  successText: {
    fontSize: 18,
    fontWeight: '600',
    color: COLORS.success,
    textAlign: 'center',
    marginTop: 16,
  },
  errorText: {
    fontSize: 14,
    color: COLORS.error,
    textAlign: 'center',
    marginBottom: 16,
  },
  retryButton: {
    marginTop: 8,
  },
  cancelButton: {
    marginTop: 8,
  },
  backButton: {
    marginTop: 16,
  },
});
