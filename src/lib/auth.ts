export type BackendUser = {
  id: string;
  email: string;
  full_name: string;
  is_active: boolean;
};

export function displayName(user: BackendUser | null) {
  return user?.full_name || "User";
}
