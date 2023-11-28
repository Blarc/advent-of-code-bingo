import {HttpClient, HttpHeaders} from '@angular/common/http';
import {Injectable, inject} from '@angular/core';

import {Observable} from 'rxjs';

import {CookiesService} from '../../services/cookies.service';
import {BingoCardDto} from '../model/bingoCardDto';
import {UserDto} from '../model/userDto';

@Injectable({
    providedIn: 'root'
})
export class BingoApiService {
    private http = inject(HttpClient);
    private cookiesService = inject(CookiesService);

    private readonly baseUrl = '/api/v1';

    private getAuthHeaders(): HttpHeaders | undefined {
        const token = this.cookiesService.getCookie('token');

        if (!token) {
            return;
        }

        return new HttpHeaders({
            'Content-Type': 'application/json',
            Authorization: `Bearer ${decodeURIComponent(token)}`
        });
    }

    public getAllBingoCards(): Observable<BingoCardDto[]> {
        const headers = this.getAuthHeaders();
        const httpOptions = headers ? {headers} : undefined;
        return this.http.get<BingoCardDto[]>(`${this.baseUrl}/bingoCards`, httpOptions);
    }

    public getUserInfo(): Observable<UserDto> {
        const headers = this.getAuthHeaders();
        return this.http.get<UserDto>(`${this.baseUrl}/me`, {headers: headers});
    }

    public clickBingoCard(uuid: string): Observable<UserDto> {
        const headers = this.getAuthHeaders();
        return this.http.post<UserDto>(`${this.baseUrl}/me/bingoCard/${uuid}`, null, {headers: headers});
    }
}
