export interface Service {
  id: string;
  title: string;
  url: string;
  description: string;
  icon: string;
  status_check: boolean;
  status_check_url: string | null;
  sort_order: number;
  created_at: string;
  updated_at: string;
}

export interface ServiceWithSections extends Service {
  section_ids: string[];
}

export interface Section {
  id: string;
  name: string;
  icon: string;
  cols: number;
  collapsed: boolean;
  sort_order: number;
  section_type: string;
  created_at: string;
  updated_at: string;
}

export interface NestedSection extends Section {
  services: Service[];
}
