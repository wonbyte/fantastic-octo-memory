import React from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  ActivityIndicator,
} from 'react-native';
import { useLocalSearchParams, router } from 'expo-router';
import { useProject } from '../../../src/hooks/useProjects';
import { useBlueprints } from '../../../src/hooks/useBlueprints';
import { Card } from '../../../src/components/ui/Card';
import { Button } from '../../../src/components/ui/Button';
import { Loading } from '../../../src/components/ui/Loading';
import { ErrorState } from '../../../src/components/ui/ErrorState';
import { COLORS, STATUS_COLORS } from '../../../src/utils/constants';
import { Blueprint } from '../../../src/types';

export default function ProjectDetailScreen() {
  const { id: projectId } = useLocalSearchParams<{ id: string }>();
  const { data: project, isLoading: projectLoading, error: projectError, refetch: refetchProject } = useProject(projectId);
  const { data: blueprints, isLoading: blueprintsLoading, error: blueprintsError, refetch: refetchBlueprints } = useBlueprints(projectId);

  const renderBlueprint = (blueprint: Blueprint) => (
    <TouchableOpacity
      key={blueprint.id}
      style={styles.blueprintItem}
      activeOpacity={0.7}
      onPress={() => router.push(`/(main)/projects/${projectId}/blueprints/${blueprint.id}`)}
    >
      <View style={styles.blueprintInfo}>
        <Text style={styles.blueprintName}>{blueprint.filename}</Text>
        <View style={styles.statusRow}>
          <View style={[styles.statusDot, { backgroundColor: STATUS_COLORS[blueprint.upload_status] }]} />
          <Text style={styles.statusLabel}>Upload: {blueprint.upload_status}</Text>
        </View>
        <View style={styles.statusRow}>
          <View style={[styles.statusDot, { backgroundColor: STATUS_COLORS[blueprint.analysis_status] }]} />
          <Text style={styles.statusLabel}>Analysis: {blueprint.analysis_status}</Text>
        </View>
      </View>
    </TouchableOpacity>
  );

  if (projectLoading) {
    return <Loading message="Loading project..." />;
  }

  if (projectError || !project) {
    return (
      <ErrorState
        message="Failed to load project"
        onRetry={refetchProject}
      />
    );
  }

  return (
    <ScrollView style={styles.container}>
      <Card style={styles.projectCard}>
        <View style={styles.header}>
          <View style={styles.headerContent}>
            <Text style={styles.projectName}>{project.name}</Text>
            <View style={[styles.statusBadge, { backgroundColor: STATUS_COLORS[project.status] }]}>
              <Text style={styles.statusText}>{project.status}</Text>
            </View>
          </View>
        </View>
        
        {project.description && (
          <Text style={styles.projectDescription}>{project.description}</Text>
        )}
        
        <View style={styles.infoRow}>
          <Text style={styles.infoLabel}>Created:</Text>
          <Text style={styles.infoValue}>
            {new Date(project.created_at).toLocaleDateString()}
          </Text>
        </View>
        
        <View style={styles.infoRow}>
          <Text style={styles.infoLabel}>Updated:</Text>
          <Text style={styles.infoValue}>
            {new Date(project.updated_at).toLocaleDateString()}
          </Text>
        </View>
      </Card>

      <Card style={styles.blueprintsCard}>
        <View style={styles.blueprintsHeader}>
          <Text style={styles.sectionTitle}>Blueprints</Text>
          <Button
            title="Upload"
            onPress={() => router.push(`/(main)/projects/${projectId}/blueprints/upload`)}
            style={styles.uploadButton}
          />
        </View>

        {blueprintsLoading && (
          <View style={styles.loadingContainer}>
            <ActivityIndicator color={COLORS.primary} />
            <Text style={styles.loadingText}>Loading blueprints...</Text>
          </View>
        )}

        {blueprintsError && (
          <View style={styles.errorContainer}>
            <Text style={styles.errorText}>Failed to load blueprints</Text>
            <Button
              title="Retry"
              onPress={refetchBlueprints}
              variant="secondary"
              style={styles.retryButton}
            />
          </View>
        )}

        {!blueprintsLoading && !blueprintsError && blueprints && blueprints.length === 0 && (
          <View style={styles.emptyContainer}>
            <Text style={styles.emptyEmoji}>📄</Text>
            <Text style={styles.emptyText}>No blueprints yet</Text>
            <Text style={styles.emptySubtext}>Upload a blueprint to get started</Text>
          </View>
        )}

        {!blueprintsLoading && !blueprintsError && blueprints && blueprints.length > 0 && (
          <View style={styles.blueprintsList}>
            {blueprints.map(renderBlueprint)}
          </View>
        )}
      </Card>

      <Button
        title="Back to Projects"
        onPress={() => router.back()}
        variant="secondary"
        style={styles.backButton}
      />
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: COLORS.background.secondary,
  },
  projectCard: {
    margin: 16,
  },
  header: {
    marginBottom: 16,
  },
  headerContent: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  projectName: {
    fontSize: 24,
    fontWeight: 'bold',
    color: COLORS.text.primary,
    flex: 1,
  },
  statusBadge: {
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: 12,
    marginLeft: 12,
  },
  statusText: {
    fontSize: 12,
    fontWeight: '600',
    color: '#FFFFFF',
  },
  projectDescription: {
    fontSize: 16,
    color: COLORS.text.secondary,
    marginBottom: 16,
    lineHeight: 24,
  },
  infoRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: 8,
  },
  infoLabel: {
    fontSize: 14,
    color: COLORS.text.secondary,
    fontWeight: '600',
  },
  infoValue: {
    fontSize: 14,
    color: COLORS.text.primary,
  },
  blueprintsCard: {
    margin: 16,
    marginTop: 0,
  },
  blueprintsHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 16,
  },
  sectionTitle: {
    fontSize: 20,
    fontWeight: '600',
    color: COLORS.text.primary,
  },
  uploadButton: {
    minWidth: 100,
  },
  loadingContainer: {
    alignItems: 'center',
    padding: 24,
  },
  loadingText: {
    marginTop: 8,
    fontSize: 14,
    color: COLORS.text.secondary,
  },
  errorContainer: {
    alignItems: 'center',
    padding: 24,
  },
  errorText: {
    fontSize: 14,
    color: COLORS.error,
    marginBottom: 16,
  },
  retryButton: {
    minWidth: 100,
  },
  emptyContainer: {
    alignItems: 'center',
    padding: 32,
  },
  emptyEmoji: {
    fontSize: 48,
    marginBottom: 12,
  },
  emptyText: {
    fontSize: 16,
    fontWeight: '600',
    color: COLORS.text.primary,
    marginBottom: 4,
  },
  emptySubtext: {
    fontSize: 14,
    color: COLORS.text.secondary,
  },
  blueprintsList: {
    marginTop: 8,
  },
  blueprintItem: {
    borderTopWidth: 1,
    borderTopColor: COLORS.border,
    paddingVertical: 12,
  },
  blueprintInfo: {
    flex: 1,
  },
  blueprintName: {
    fontSize: 16,
    fontWeight: '600',
    color: COLORS.text.primary,
    marginBottom: 8,
  },
  statusRow: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 4,
  },
  statusDot: {
    width: 8,
    height: 8,
    borderRadius: 4,
    marginRight: 8,
  },
  statusLabel: {
    fontSize: 14,
    color: COLORS.text.secondary,
  },
  backButton: {
    margin: 16,
    marginTop: 8,
  },
});
