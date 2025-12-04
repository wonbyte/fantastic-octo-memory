import React from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  ActivityIndicator,
} from 'react-native';
import { useLocalSearchParams, router } from 'expo-router';
import { useJob } from '../../../../../../src/hooks/useJobs';
import { Card } from '../../../../../../src/components/ui/Card';
import { Button } from '../../../../../../src/components/ui/Button';
import { Loading } from '../../../../../../src/components/ui/Loading';
import { ErrorState } from '../../../../../../src/components/ui/ErrorState';
import { COLORS, STATUS_COLORS } from '../../../../../../src/utils/constants';

export default function AnalysisStatusScreen() {
  const { jobId, id: projectId, blueprintId } = useLocalSearchParams<{ 
    jobId: string; 
    id: string; 
    blueprintId: string;
  }>();
  const { data: job, isLoading, error, refetch } = useJob(jobId, !!jobId);

  const getStatusIcon = () => {
    switch (job?.status) {
      case 'queued':
        return '⏳';
      case 'processing':
        return '🔄';
      case 'completed':
        return '✓';
      case 'failed':
        return '✗';
      default:
        return '•';
    }
  };

  const getStatusMessage = () => {
    switch (job?.status) {
      case 'queued':
        return 'Waiting in queue...';
      case 'processing':
        return 'Analyzing blueprint...';
      case 'completed':
        return 'Analysis complete!';
      case 'failed':
        return 'Analysis failed';
      default:
        return 'Unknown status';
    }
  };

  const getElapsedTime = () => {
    if (!job?.created_at) return null;
    
    const start = new Date(job.created_at).getTime();
    // For completed jobs, use completed_at time
    // For in-progress jobs, we can't use Date.now() during render
    // So we return null to avoid showing misleading time
    if (!job.completed_at) {
      return null;
    }
    
    const end = new Date(job.completed_at).getTime();
    const elapsed = Math.floor((end - start) / 1000);
    
    if (elapsed < 60) return `${elapsed}s`;
    const minutes = Math.floor(elapsed / 60);
    const seconds = elapsed % 60;
    return `${minutes}m ${seconds}s`;
  };

  if (isLoading) {
    return <Loading message="Loading job status..." />;
  }

  if (error || !job) {
    return (
      <ErrorState
        message="Failed to load job status"
        onRetry={refetch}
      />
    );
  }

  const isProcessing = job.status === 'queued' || job.status === 'processing';

  return (
    <ScrollView style={styles.container}>
      <Card style={styles.card}>
        <View style={styles.statusHeader}>
          <Text style={styles.statusIcon}>{getStatusIcon()}</Text>
          <Text style={styles.statusTitle}>{getStatusMessage()}</Text>
        </View>

        {isProcessing && (
          <View style={styles.processingContainer}>
            <ActivityIndicator size="large" color={COLORS.primary} />
            <Text style={styles.processingText}>
              Please wait while we analyze your blueprint...
            </Text>
          </View>
        )}

        <View style={styles.infoSection}>
          <View style={styles.infoRow}>
            <Text style={styles.infoLabel}>Job ID:</Text>
            <Text style={styles.infoValue}>{job.id}</Text>
          </View>

          <View style={styles.infoRow}>
            <Text style={styles.infoLabel}>Status:</Text>
            <View style={styles.statusContainer}>
              <View
                style={[
                  styles.statusDot,
                  { backgroundColor: STATUS_COLORS[job.status] },
                ]}
              />
              <Text style={styles.infoValue}>{job.status}</Text>
            </View>
          </View>

          {job.progress !== undefined && job.progress > 0 && (
            <View style={styles.progressSection}>
              <View style={styles.infoRow}>
                <Text style={styles.infoLabel}>Progress:</Text>
                <Text style={styles.infoValue}>{job.progress}%</Text>
              </View>
              <View style={styles.progressBar}>
                <View
                  style={[
                    styles.progressFill,
                    { width: `${job.progress}%` },
                  ]}
                />
              </View>
            </View>
          )}

          {getElapsedTime() && (
            <View style={styles.infoRow}>
              <Text style={styles.infoLabel}>
                {job.status === 'completed' ? 'Duration:' : 'Elapsed Time:'}
              </Text>
              <Text style={styles.infoValue}>{getElapsedTime()}</Text>
            </View>
          )}

          {job.created_at && (
            <View style={styles.infoRow}>
              <Text style={styles.infoLabel}>Started:</Text>
              <Text style={styles.infoValue}>
                {new Date(job.created_at).toLocaleString()}
              </Text>
            </View>
          )}

          {job.completed_at && (
            <View style={styles.infoRow}>
              <Text style={styles.infoLabel}>Completed:</Text>
              <Text style={styles.infoValue}>
                {new Date(job.completed_at).toLocaleString()}
              </Text>
            </View>
          )}
        </View>

        {job.error_message && (
          <Card style={styles.errorCard}>
            <Text style={styles.errorTitle}>Error Details</Text>
            <Text style={styles.errorMessage}>{job.error_message}</Text>
          </Card>
        )}

        {job.status === 'completed' && job.result && (
          <>
            <Card style={styles.resultsCard}>
              <Text style={styles.resultsTitle}>Analysis Results</Text>
              
              {job.result.summary && (
                <View style={styles.summarySection}>
                  <Text style={styles.summaryTitle}>Summary</Text>
                  <View style={styles.summaryGrid}>
                    <View style={styles.summaryItem}>
                      <Text style={styles.summaryValue}>
                        {job.result.summary.total_rooms || 0}
                      </Text>
                      <Text style={styles.summaryLabel}>Rooms</Text>
                    </View>
                    <View style={styles.summaryItem}>
                      <Text style={styles.summaryValue}>
                        {job.result.summary.total_openings || 0}
                      </Text>
                      <Text style={styles.summaryLabel}>Openings</Text>
                    </View>
                    <View style={styles.summaryItem}>
                      <Text style={styles.summaryValue}>
                        {job.result.summary.total_fixtures || 0}
                      </Text>
                      <Text style={styles.summaryLabel}>Fixtures</Text>
                    </View>
                    {job.result.summary.total_area && (
                      <View style={styles.summaryItem}>
                        <Text style={styles.summaryValue}>
                          {job.result.summary.total_area.toFixed(0)}
                        </Text>
                        <Text style={styles.summaryLabel}>sq ft</Text>
                      </View>
                    )}
                  </View>
                </View>
              )}

              {/* Rooms Section */}
              {job.result.rooms && job.result.rooms.length > 0 && (
                <View style={styles.detailSection}>
                  <Text style={styles.detailTitle}>Rooms</Text>
                  {job.result.rooms.map((room, index) => (
                    <View key={index} style={styles.detailItem}>
                      <Text style={styles.detailName}>
                        {room.name}
                        {room.room_type && (
                          <Text style={styles.detailType}> ({room.room_type})</Text>
                        )}
                      </Text>
                      <Text style={styles.detailInfo}>
                        {room.dimensions} • {room.area.toFixed(0)} sq ft
                      </Text>
                    </View>
                  ))}
                </View>
              )}

              {/* Openings Section */}
              {job.result.openings && job.result.openings.length > 0 && (
                <View style={styles.detailSection}>
                  <Text style={styles.detailTitle}>Openings (Doors & Windows)</Text>
                  {job.result.openings.map((opening, index) => (
                    <View key={index} style={styles.detailItem}>
                      <Text style={styles.detailName}>
                        {opening.opening_type}
                      </Text>
                      <Text style={styles.detailInfo}>
                        {opening.count} × {opening.size}
                        {opening.details && ` • ${opening.details}`}
                      </Text>
                    </View>
                  ))}
                </View>
              )}

              {/* Fixtures Section */}
              {job.result.fixtures && job.result.fixtures.length > 0 && (
                <View style={styles.detailSection}>
                  <Text style={styles.detailTitle}>Fixtures</Text>
                  {job.result.fixtures.map((fixture, index) => (
                    <View key={index} style={styles.detailItem}>
                      <Text style={styles.detailName}>
                        {fixture.fixture_type}
                        {fixture.category && (
                          <Text style={styles.detailType}> ({fixture.category})</Text>
                        )}
                      </Text>
                      <Text style={styles.detailInfo}>
                        Qty: {fixture.count}
                        {fixture.details && ` • ${fixture.details}`}
                      </Text>
                    </View>
                  ))}
                </View>
              )}

              {/* Materials Section */}
              {job.result.materials && job.result.materials.length > 0 && (
                <View style={styles.detailSection}>
                  <Text style={styles.detailTitle}>Materials</Text>
                  {job.result.materials.map((material, index) => (
                    <View key={index} style={styles.detailItem}>
                      <Text style={styles.detailName}>
                        {material.material_name}
                      </Text>
                      <Text style={styles.detailInfo}>
                        {material.quantity} {material.unit}
                        {material.specifications && ` • ${material.specifications}`}
                      </Text>
                    </View>
                  ))}
                </View>
              )}
            </Card>

            {/* Generate Bid Button */}
            <View style={styles.bidSection}>
              <Button
                title="Generate Bid"
                onPress={() => {
                  // Navigate to bid generation using params
                  router.push(`/projects/${projectId}/generate-bid?blueprintId=${blueprintId}`);
                }}
                style={styles.bidButton}
              />
            </View>
          </>
        )}

        <View style={styles.actionSection}>
          {job.status === 'failed' && (
            <Button
              title="Retry Analysis"
              onPress={() => {
                // Trigger new analysis
                router.back();
              }}
              style={styles.retryButton}
            />
          )}

          <Button
            title="Back to Blueprint"
            onPress={() => router.back()}
            variant="secondary"
            style={styles.backButton}
          />
        </View>
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
  statusHeader: {
    alignItems: 'center',
    marginBottom: 24,
  },
  statusIcon: {
    fontSize: 64,
    marginBottom: 16,
  },
  statusTitle: {
    fontSize: 24,
    fontWeight: 'bold',
    color: COLORS.text.primary,
  },
  processingContainer: {
    alignItems: 'center',
    paddingVertical: 24,
    marginBottom: 24,
  },
  processingText: {
    marginTop: 16,
    fontSize: 14,
    color: COLORS.text.secondary,
    textAlign: 'center',
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
  progressSection: {
    marginTop: 8,
  },
  progressBar: {
    height: 8,
    backgroundColor: COLORS.border,
    borderRadius: 4,
    overflow: 'hidden',
    marginTop: 8,
  },
  progressFill: {
    height: '100%',
    backgroundColor: COLORS.primary,
    borderRadius: 4,
  },
  errorCard: {
    backgroundColor: '#FEE2E2',
    borderColor: COLORS.error,
    borderWidth: 1,
    marginBottom: 24,
  },
  errorTitle: {
    fontSize: 16,
    fontWeight: '600',
    color: COLORS.error,
    marginBottom: 8,
  },
  errorMessage: {
    fontSize: 14,
    color: '#991B1B',
  },
  resultsCard: {
    backgroundColor: '#ECFDF5',
    borderColor: COLORS.success,
    borderWidth: 1,
    marginBottom: 24,
  },
  resultsTitle: {
    fontSize: 20,
    fontWeight: '600',
    color: COLORS.text.primary,
    marginBottom: 16,
  },
  summarySection: {
    marginBottom: 16,
  },
  summaryTitle: {
    fontSize: 16,
    fontWeight: '600',
    color: COLORS.text.secondary,
    marginBottom: 12,
  },
  summaryGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    justifyContent: 'space-between',
  },
  summaryItem: {
    width: '48%',
    alignItems: 'center',
    padding: 16,
    backgroundColor: COLORS.background.primary,
    borderRadius: 8,
    marginBottom: 8,
  },
  summaryValue: {
    fontSize: 32,
    fontWeight: 'bold',
    color: COLORS.primary,
    marginBottom: 4,
  },
  summaryLabel: {
    fontSize: 14,
    color: COLORS.text.secondary,
  },
  actionSection: {
    marginTop: 16,
  },
  retryButton: {
    marginBottom: 12,
  },
  backButton: {
    marginTop: 8,
  },
  detailSection: {
    marginTop: 24,
    paddingTop: 16,
    borderTopWidth: 1,
    borderTopColor: COLORS.border,
  },
  detailTitle: {
    fontSize: 18,
    fontWeight: '600',
    color: COLORS.text.primary,
    marginBottom: 12,
  },
  detailItem: {
    paddingVertical: 12,
    borderBottomWidth: 1,
    borderBottomColor: COLORS.background.secondary,
  },
  detailName: {
    fontSize: 16,
    fontWeight: '500',
    color: COLORS.text.primary,
    marginBottom: 4,
  },
  detailType: {
    fontSize: 14,
    fontWeight: '400',
    color: COLORS.text.secondary,
  },
  detailInfo: {
    fontSize: 14,
    color: COLORS.text.secondary,
  },
  bidSection: {
    marginTop: 16,
    marginBottom: 8,
  },
  bidButton: {
    backgroundColor: COLORS.success,
  },
});
