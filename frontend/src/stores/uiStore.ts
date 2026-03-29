import { create } from 'zustand'

const STORAGE_KEY = 'dashboard-ui-collapsed-sections'

type PersistedUIState = {
  collapsedSections: string[]
  knownSections: string[]
}

export type SectionDefault = {
  id: string
  collapsed: boolean
}

interface UIStoreState {
  collapsedSections: Set<string>
  knownSections: Set<string>
  toggleSection: (id: string) => void
  isCollapsed: (id: string) => boolean
  initializeSections: (sections: SectionDefault[]) => void
}

const loadPersistedState = (): PersistedUIState => {
  if (typeof window === 'undefined') {
    return { collapsedSections: [], knownSections: [] }
  }

  const raw = window.localStorage.getItem(STORAGE_KEY)
  if (!raw) {
    return { collapsedSections: [], knownSections: [] }
  }

  try {
    const parsed = JSON.parse(raw) as Partial<PersistedUIState>
    return {
      collapsedSections: Array.isArray(parsed.collapsedSections) ? parsed.collapsedSections : [],
      knownSections: Array.isArray(parsed.knownSections) ? parsed.knownSections : [],
    }
  } catch {
    return { collapsedSections: [], knownSections: [] }
  }
}

const persistState = (collapsedSections: Set<string>, knownSections: Set<string>) => {
  if (typeof window === 'undefined') {
    return
  }

  window.localStorage.setItem(
    STORAGE_KEY,
    JSON.stringify({
      collapsedSections: [...collapsedSections],
      knownSections: [...knownSections],
    } satisfies PersistedUIState),
  )
}

const persisted = loadPersistedState()

export const useUIStore = create<UIStoreState>((set, get) => ({
  collapsedSections: new Set(persisted.collapsedSections),
  knownSections: new Set(persisted.knownSections),
  toggleSection: (id) =>
    set((state) => {
      const collapsedSections = new Set(state.collapsedSections)
      const knownSections = new Set(state.knownSections)
      knownSections.add(id)

      if (collapsedSections.has(id)) {
        collapsedSections.delete(id)
      } else {
        collapsedSections.add(id)
      }

      persistState(collapsedSections, knownSections)
      return { collapsedSections, knownSections }
    }),
  isCollapsed: (id) => get().collapsedSections.has(id),
  initializeSections: (sections) =>
    set((state) => {
      const collapsedSections = new Set(state.collapsedSections)
      const knownSections = new Set(state.knownSections)
      let changed = false

      for (const section of sections) {
        if (knownSections.has(section.id)) {
          continue
        }
        knownSections.add(section.id)
        changed = true
        if (section.collapsed) {
          collapsedSections.add(section.id)
        }
      }

      if (!changed) {
        return state
      }

      persistState(collapsedSections, knownSections)
      return { collapsedSections, knownSections }
    }),
}))
