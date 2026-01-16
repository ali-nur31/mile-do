export const combineToBackend = (dateStr: string, timeStr?: string): string => {
  if (!dateStr || dateStr.trim() === '') return "";
  
  let finalTime = "09:00:00";
  
  if (timeStr && timeStr.trim() !== '') {
    const timeParts = timeStr.split(':');
    const hours = (timeParts[0] || "09").padStart(2, '0');
    const minutes = (timeParts[1] || "00").padStart(2, '0');
    const seconds = (timeParts[2] || "00").padStart(2, '0');
    finalTime = `${hours}:${minutes}:${seconds}`;
  }
  
  return `${dateStr.trim()} ${finalTime}`;
};

export const extractDateStr = (isoString?: string): string => {
  if (!isoString) return "";
  
  const datePart = isoString.split(' ')[0];
  
  if (datePart.startsWith('0001')) return "";
  
  return datePart;
};

export const extractTimeStr = (isoString?: string): string => {
  if (!isoString) return "";
  
  const parts = isoString.split(' ');
  if (parts.length < 2) return "";
  
  const timePart = parts[1];
  const match = timePart.match(/(\d{2}):(\d{2}):(\d{2})/);
  
  if (match) {
    return `${match[1]}:${match[2]}`;
  }
  return "";
};

export const formatDisplayDate = (isoString?: string): string => {
  if (!isoString || isoString.startsWith('0001')) return "";
  const date = new Date(isoString);
  if (isNaN(date.getTime())) return "";
  return date.toLocaleDateString(undefined, { month: 'short', day: 'numeric' });
};

export const getLocalTodayStr = (): string => {
  const d = new Date();
  const year = d.getFullYear();
  const month = String(d.getMonth() + 1).padStart(2, '0');
  const day = String(d.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
};

export const safeDate = (dateStr?: string): Date | null => {
  if (!dateStr || dateStr.startsWith('0001')) return null;
  const d = new Date(dateStr);
  return isNaN(d.getTime()) ? null : d;
};

export const getLocalISOString = (date: Date): string => {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
};

export const extractTimeFromBackend = extractTimeStr;
export const getTodayStr = getLocalTodayStr;