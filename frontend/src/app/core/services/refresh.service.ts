import {Injectable} from '@angular/core';

import {BehaviorSubject, Observable} from 'rxjs';

@Injectable({providedIn: 'root'})
export class RefreshService {
    private refreshBingoCards = new BehaviorSubject(false);

    shouldRefreshBingoCards(): void {
        this.refreshBingoCards.next(true);
    }

    onRefreshBingoCards(): Observable<boolean> {
        return this.refreshBingoCards.asObservable();
    }
}
