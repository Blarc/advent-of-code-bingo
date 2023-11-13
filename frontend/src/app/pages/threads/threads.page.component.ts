import {Component} from '@angular/core';

@Component({
    standalone: true,
    selector: 'app-threads',
    templateUrl: 'threads.page.component.html',
    styleUrls: ['threads.page.component.scss']
})
export class ThreadsPageComponent {
    public voteCounter = 0;

    upVote() {
        this.voteCounter++;
    }

    downVote() {
        this.voteCounter--;
    }
}
