export const combineToBackend = (dateStr: string, timeStr?: string): string | undefined => {
  if (!dateStr || dateStr.trim() === '') return undefined;
  
  let finalTime = "09:00:00";
  
  if (timeStr && timeStr.trim() !== '') {
    const parts = timeStr.split(':');
    const h = (parts[0] || "09").padStart(2, '0');
    const m = (parts[1] || "00").padStart(2, '0');
    const s = (parts[2] || "00").padStart(2, '0');
    finalTime = `${h}:${m}:${s}`;
  }
  
  return `${dateStr.trim()} ${finalTime}`;
};

export const extractDateStr = (isoString?: string): string => {
  if (!isoString) return "";
  const datePart = isoString.split(/[ T]/)[0];
  if (datePart.startsWith('0001')) return "";
  return datePart;
};

export const extractTimeStr = (isoString?: string): string => {
  if (!isoString) return "";
  const parts = isoString.split(/[ T]/);
  if (parts.length < 2) return "";
  
  const timePart = parts[1];
  const match = timePart.match(/(\d{2}):(\d{2}):(\d{2})/);
  if (match) return `${match[1]}:${match[2]}`;
  return "";
};

export const extractTimeFromBackend = extractTimeStr;

export const getLocalTodayStr = (): string => {
  const d = new Date();
  const year = d.getFullYear();
  const month = String(d.getMonth() + 1).padStart(2, '0');
  const day = String(d.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
};

export const getTodayStr = getLocalTodayStr;

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

export const addMinutesToTime = (timeStr: string, minutesToAdd: number): string => {
  if (!timeStr) return "09:00";
  const [hours, minutes] = timeStr.split(':').map(Number);
  const date = new Date();
  date.setHours(hours || 0, minutes || 0, 0, 0);
  date.setMinutes(date.getMinutes() + minutesToAdd);
  const h = String(date.getHours()).padStart(2, '0');
  const m = String(date.getMinutes()).padStart(2, '0');
  return `${h}:${m}`;
};

export const calculateDuration = (startTimeStr: string, endTimeStr: string): number => {
  if (!startTimeStr || !endTimeStr) return 0;
  const [startH, startM] = startTimeStr.split(':').map(Number);
  const [endH, endM] = endTimeStr.split(':').map(Number);
  const startTotal = (startH || 0) * 60 + (startM || 0);
  const endTotal = (endH || 0) * 60 + (endM || 0);
  let diff = endTotal - startTotal;
  if (diff < 0) diff += 24 * 60;
  return diff;
};