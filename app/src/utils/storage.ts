// Storage helper for non-secure data
// For web, uses localStorage; for native, could use AsyncStorage
// For now, using a simple in-memory implementation with localStorage fallback

class Storage {
  private memoryStore: Map<string, string> = new Map();

  async setItem(key: string, value: string): Promise<void> {
    try {
      this.memoryStore.set(key, value);
      if (typeof window !== 'undefined' && window.localStorage) {
        window.localStorage.setItem(key, value);
      }
    } catch (error) {
      console.error('Error saving to storage:', error);
      throw error;
    }
  }

  async getItem(key: string): Promise<string | null> {
    try {
      // Try memory first
      if (this.memoryStore.has(key)) {
        return this.memoryStore.get(key) || null;
      }
      // Fallback to localStorage on web
      if (typeof window !== 'undefined' && window.localStorage) {
        const value = window.localStorage.getItem(key);
        if (value) {
          this.memoryStore.set(key, value);
        }
        return value;
      }
      return null;
    } catch (error) {
      console.error('Error reading from storage:', error);
      return null;
    }
  }

  async removeItem(key: string): Promise<void> {
    try {
      this.memoryStore.delete(key);
      if (typeof window !== 'undefined' && window.localStorage) {
        window.localStorage.removeItem(key);
      }
    } catch (error) {
      console.error('Error removing from storage:', error);
      throw error;
    }
  }

  async clear(): Promise<void> {
    try {
      this.memoryStore.clear();
      if (typeof window !== 'undefined' && window.localStorage) {
        window.localStorage.clear();
      }
    } catch (error) {
      console.error('Error clearing storage:', error);
      throw error;
    }
  }
}

export const storage = new Storage();
