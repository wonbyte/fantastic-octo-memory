import React from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  Alert,
} from 'react-native';
import { useLocalSearchParams, router } from 'expo-router';
import { useBlueprint, useTriggerAnalysis } from '../../../../../src/hooks/useBlueprints';
import { useJobsByBlueprint } from '../../../../../src/hooks/useJobs';
import { Card } from '../../../../../src/components/ui/Card';
import { Button } from '../../../../../src/components/ui/Button';
import { Loading } from '../../../../../src/components/ui/Loading';
import { ErrorState } from '../../../../../src/components/ui/ErrorState';
import { COLORS, STATUS_COLORS } from '../../../../../src/utils/constants';

export default function BlueprintDetailScreen() {
  const { blueprintId, id: projectId } = useLocalSearchParams<{ blueprintId: string; id: string }>();
  const { data: blueprint, isLoading, error, refetch } = useBlueprint(blueprintId);
  const { data: jobs } = useJobsByBlueprint(blueprintId);
  const triggerAnalysis = useTriggerAnalysis();

  const handleAnalyze = async () => {
    if (!blueprintId || !projectId) return;

    try {
      const result = await triggerAnalysis.mutateAsync(blueprintId);
      Alert.alert(
        'Analysis Started',
        'Blueprint analysis has been queued.',
        [
          {
            text: 'View Status',
            onPress: () => router.push(`/(main)/projects/${projectId}/blueprints/${blueprintId}/analysis?jobId=${result.job_id}`),
          },
          {
            text: 'OK',
            style: 'cancel',
          },
        ]
      );
    } catch (err) {
      console.error('Failed to trigger analysis:', err);
      Alert.alert(
        'Error',
        'Failed to start analysis. Please try again.',
        [{ text: 'OK' }]
      );
    }
  };

  if (isLoading) {
    return <Loading message="Loading blueprint..." />;
  }

  if (error || !blueprint) {
    return (
      <ErrorState
        message="Failed to load blueprint"
        onRetry={refetch}
      />
    );
  }

  const canAnalyze = blueprint.upload_status === 'uploaded' && 
                     blueprint.analysis_status !== 'processing' &&
                     blueprint.analysis_status !== 'queued';

  const latestJob = jobs && jobs.length > 0 ? jobs[0] : null;

  return (
    <ScrollView style={styles.container}>
      <Card style={styles.card}>
        <Text style={styles.title}>{blueprint.filename}</Text>

        <View style={styles.infoSection}>
          <View style={styles.infoRow}>
            <Text style={styles.infoLabel}>File Size:</Text>
            <Text style={styles.infoValue}>
              {(blueprint.file_size / (1024 * 1024)).toFixed(2)} MB
            </Text>
          </View>

          <View style={styles.infoRow}>
            <Text style={styles.infoLabel}>Content Type:</Text>
            <Text style={styles.infoValue}>{blueprint.content_type}</Text>
          </View>

          <View style={styles.infoRow}>
            <Text style={styles.infoLabel}>Upload Status:</Text>
            <View style={styles.statusContainer}>
              <View
                style={[
                  styles.statusDot,
                  { backgroundColor: STATUS_COLORS[blueprint.upload_status] },
                ]}
              />
              <Text style={styles.infoValue}>{blueprint.upload_status}</Text>
            </View>
          </View>

          <View style={styles.infoRow}>
            <Text style={styles.infoLabel}>Analysis Status:</Text>
            <View style={styles.statusContainer}>
              <View
                style={[
                  styles.statusDot,
                  { backgroundColor: STATUS_COLORS[blueprint.analysis_status] },
                ]}
              />
              <Text style={styles.infoValue}>{blueprint.analysis_status}</Text>
            </View>
          </View>

          <View style={styles.infoRow}>
            <Text style={styles.infoLabel}>Uploaded:</Text>
            <Text style={styles.infoValue}>
              {new Date(blueprint.created_at).toLocaleString()}
            </Text>
          </View>
        </View>

        {blueprint.upload_status === 'uploaded' && (
          <View style={styles.actionSection}>
            <Button
              title={blueprint.analysis_status === 'completed' ? 'Analyze Again' : 'Analyze Blueprint'}
              onPress={handleAnalyze}
              disabled={!canAnalyze || triggerAnalysis.isPending}
              loading={triggerAnalysis.isPending}
              style={styles.analyzeButton}
            />

            {!canAnalyze && blueprint.analysis_status === 'processing' && (
              <Text style={styles.helpText}>
                Analysis is currently in progress
              </Text>
            )}
            {!canAnalyze && blueprint.analysis_status === 'queued' && (
              <Text style={styles.helpText}>
                Analysis is queued and will start soon
              </Text>
            )}
          </View>
        )}

        {latestJob && (
          <Card style={styles.jobCard}>
            <Text style={styles.jobTitle}>Latest Analysis Job</Text>
            <View style={styles.infoRow}>
              <Text style={styles.infoLabel}>Status:</Text>
              <View style={styles.statusContainer}>
                <View
                  style={[
                    styles.statusDot,
                    { backgroundColor: STATUS_COLORS[latestJob.status] },
                  ]}
                />
                <Text style={styles.infoValue}>{latestJob.status}</Text>
              </View>
            </View>
            {latestJob.progress !== undefined && (
              <View style={styles.infoRow}>
                <Text style={styles.infoLabel}>Progress:</Text>
                <Text style={styles.infoValue}>{latestJob.progress}%</Text>
              </View>
            )}
            <Button
              title="View Analysis Status"
              onPress={() => router.push(`/(main)/projects/${projectId}/blueprints/${blueprintId}/analysis?jobId=${latestJob.id}`)}
              variant="secondary"
              style={styles.viewJobButton}
            />
          </Card>
        )}

        <Button
          title="Back to Project"
          onPress={() => router.back()}
          variant="secondary"
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
    marginBottom: 24,
  },
  infoSection: {
    marginBottom: 24,
  },
  infoRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 12,
  },
  infoLabel: {
    fontSize: 14,
    fontWeight: '600',
    color: COLORS.text.secondary,
  },
  infoValue: {
    fontSize: 14,
    color: COLORS.text.primary,
  },
  statusContainer: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  statusDot: {
    width: 8,
    height: 8,
    borderRadius: 4,
    marginRight: 8,
  },
  actionSection: {
    marginBottom: 24,
  },
  analyzeButton: {
    marginBottom: 8,
  },
  helpText: {
    fontSize: 14,
    color: COLORS.text.secondary,
    textAlign: 'center',
    fontStyle: 'italic',
  },
  jobCard: {
    backgroundColor: COLORS.background.secondary,
    marginBottom: 24,
  },
  jobTitle: {
    fontSize: 18,
    fontWeight: '600',
    color: COLORS.text.primary,
    marginBottom: 16,
  },
  viewJobButton: {
    marginTop: 12,
  },
  backButton: {
    marginTop: 8,
  },
});
