export interface Service {
  id: string;
  title: string;
  url: string;
  description: string;
  icon: string;
  status_check: boolean;
  status_check_url: string | null;
  sort_order: number;
}

export interface Section {
  id: string;
  name: string;
  icon: string;
  cols: number;
  collapsed: boolean;
  sort_order: number;
  section_type: string;
  services: Service[];
}
