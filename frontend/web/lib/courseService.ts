import api from './api';

export interface Course {
  id: string;
  title: string;
  description: string;
  language: string;
  level: string;
  price: number;
  thumbnail_url: string;
  enrolled_count: number;
  rating: number;
}

export const courseService = {
  async getAllCourses(params?: { language?: string; level?: string; limit?: number }) {
    const response = await api.get('/api/v1/courses', { params });
    return response.data;
  },

  async getCourse(id: string) {
    const response = await api.get(`/api/v1/courses/${id}`);
    return response.data;
  },

  async enrollInCourse(id: string) {
    const response = await api.post(`/api/v1/courses/${id}/enroll`);
    return response.data;
  },

  async getMyCourses() {
    const response = await api.get('/api/v1/courses/my-courses');
    return response.data;
  },
};
